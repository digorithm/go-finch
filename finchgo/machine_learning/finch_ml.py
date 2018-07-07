"""
This module will contain all models
"""

import pandas as pd
import sklearn
import os
from os import listdir
from os.path import isfile, join
from sklearn import preprocessing
from sklearn.externals import joblib
from sklearn import linear_model
from sklearn.model_selection import ShuffleSplit, cross_val_score
from collections import Counter

sklearn.warnings.filterwarnings('ignore')

class FinchML:
  def __init__(self):
    self.SLI_models = None
    self.SLI_knob_models = None
    self.Encoders = None
    self.Dataset = None
    self.Scaler = None

  def train_models(self, dataset_filepath="dataset.csv", dataframe=None):
    """Train all models.
    It takes either the path or the loaded pandas dataframe"""

    if dataframe is not None:
      dataset = dataframe
    else:
      dataset = pd.read_csv(dataset_filepath)

    self.Dataset = self.preprocess_data(dataset)

    sli_models = self.train_sli_models()
    sli_knob_models = self.train_sli_to_knob_models()

    self.SLI_models = sli_models
    self.SLI_knob_models = sli_knob_models

    self.save_models(sli_models, sli_knob_models)

    return sli_knob_models, sli_models
  
  def save_models(self, sli_models={}, sli_knob_models={}):

    SLI_models_directory = 'SLI_models'
    SLI_knob_models_directory = 'SLI_knob_models'

    if not os.path.exists(SLI_models_directory):
      os.makedirs(SLI_models_directory)
    
    if not os.path.exists(SLI_knob_models_directory):
      os.makedirs(SLI_knob_models_directory)

    for model in sli_models:
      joblib.dump(sli_models[model], SLI_models_directory + '/' + model + '.pkl')
    
    for model in sli_knob_models:
      joblib.dump(sli_knob_models[model], SLI_knob_models_directory + '/' + model + '.pkl')

  def load_models(self):
    sli_models = {}
    sli_knob_models = {}
    
    sli_models_dir = 'SLI_models'
    sli_knob_models_dir = 'SLI_knob_models'

    sli_models_files = [f for f in listdir(sli_models_dir) if isfile(join(sli_models_dir, f))]

    sli_knob_models_files = [f for f in listdir(sli_knob_models_dir) if isfile(join(sli_knob_models_dir, f))]

    for file in sli_models_files:
      model_name, file_extension = os.path.splitext(file)
      model = joblib.load(sli_models_dir + '/' + file)
      sli_models[model_name] = model
    
    for file in sli_knob_models_files:
      model_name, file_extension = os.path.splitext(file)
      model = joblib.load(sli_knob_models_dir + '/' + file)
      sli_knob_models[model_name] = model

    self.SLI_models = sli_models
    self.SLI_knob_models = sli_knob_models
  
  def save_encoders(self, encoders):
    encoders_dir = 'Encoders'

    if not os.path.exists(encoders_dir):
      os.makedirs(encoders_dir)
    
    for encoder in encoders:
      joblib.dump(encoders[encoder], encoders_dir + '/' + encoder + '.pkl')
  
  def load_encoders(self):
    encoders = {}
    encoders_dir = 'Encoders'    
    encoders_files = [f for f in listdir(encoders_dir) if isfile(join(encoders_dir, f))]

    for file in encoders_files:
      encoder_name, file_extension = os.path.splitext(file)
      encoder = joblib.load(encoders_dir + '/' + file)
      encoders[encoder_name] = encoder
    
    self.Encoders = encoders
  
  def save_scaler(self, scaler):
    scaler_dir = 'Scaler'

    if not os.path.exists(scaler_dir):
      os.makedirs(scaler_dir)
    
    joblib.dump(scaler, scaler_dir + '/Scaler.pkl')
  
  def load_scaler(self):
    scaler_dir = 'Scaler'
    scaler = joblib.load(scaler_dir + '/Scaler.pkl')
    self.Scaler = scaler

  def predict_optimal_knobs(self, X, knobs, slos):
    """
    Given the system's context X, the used knobs, and the defined SLOs,
    predict the optimal set of knob values
    """
    self.load_models()
    self.load_encoders()
    self.load_scaler()

    all_predictions = {}
    for sli, slo in slos.items():
      # Our desired SLI value. We change the test sample to contain the desired SLI but keeping the current system's context
      X[sli] = slo
      X_prime = self.preprocess_data(X)
      X_prime = self.preprocess_to_predict_sli_to_knob(X_prime, sli=sli)

      
      predicted_knob_values = {}
      for knob in knobs:
        pred = self.SLI_knob_models[sli][knob].predict([X_prime.as_matrix()])
        predicted_knob_values[knob] = self.Encoders[knob].inverse_transform(pred)[0]

      all_predictions[sli] = predicted_knob_values

    final_knob_predictions = {}
    for knob in knobs:
      # Voting mechanism, it grabs the highest occurence among all predictions for this given knob
      final_knob_predictions[knob] = Counter([all_predictions[sli][knob] for sli in slos.keys()]).most_common(1)[0][0]
    
    return final_knob_predictions
    
  def predict_SLIs(self, X, slos, predicted_knobs=None):
    """
    Given system's context X, predicts all SLIs defined in the SLOs.
    If knobs are passed, it means that we are predicting SLIs for a hypothetical set of knobs, thus, we change the knob values in X to these knobs and predict how are the SLIs given the hypothetical knobs.
    """
    if predicted_knobs:
      for knob in predicted_knobs:
        X[knob] = predicted_knobs[knob]

    sli_predictions = {}
    for sli in slos.keys():
      slis_to_be_modified = [s for s in slos.keys() if s != sli]
      # Pulling the other SLIs to be equal to SLO
      for s in slis_to_be_modified:
        X[s] = slos[s]

      X_prime = self.preprocess_data(X)
      X_prime = self.preprocess_to_predict_sli(X_prime, sli)

      pred = self.SLI_models[sli].predict([X_prime])
      sli_predictions[sli] = pred[0]

    return sli_predictions

  def preprocess_data(self, original_dataset):
    remove_timestamp = False
    # make a list of desired features, drop anything other than these feats
    features = ["recipes", "houses", "schedules", "storages", "users", "CPUIdle", "CPUUsage", "HTTPRequestCount", "IOWait", "MemoryUsage"]

    all_features = original_dataset.keys().tolist()
    
    # Select which features to keep
    features_to_keep = []
    for feature in features:
      for f in all_features:
        if feature == f or feature in f or 'knob' in f:
          features_to_keep.append(f)

    features_to_keep = list(set(features_to_keep))
    features_to_drop = set(all_features) - set(features_to_keep)

    dataset = original_dataset
    
    # Remove unwanted features
    for feature in features_to_drop:
      if len(original_dataset.shape) == 1:
        dataset = dataset.drop(feature)
      else:
        dataset = dataset.drop(feature, 1)
    
    non_knob_features = [k for k in dataset.keys() if "knob" not in k]
    knob_features = [k for k in dataset.keys() if "knob" in k]

    # Unique sample
    if len(original_dataset.shape) == 1:
      scaled = self.Scaler.transform([dataset[non_knob_features]])[0]
      dataset[non_knob_features] = scaled
      
      for feat in knob_features:
        # Encode the labels for pg features
        dataset[feat] = self.Encoders[feat].transform([dataset[feat]])

    # Training phase 
    else:
      # Scale only features that are not pg variables
      self.Scaler = preprocessing.StandardScaler().fit(dataset[non_knob_features])

      dataset[non_knob_features] = self.Scaler.transform(dataset[non_knob_features])
      
      # Use encoders dict whenever you need to retrieve the real predicted label
      encoders = {}
      for feat in knob_features:
          # Encode the labels for pg features
          enc = preprocessing.LabelEncoder()
          dataset[feat] = enc.fit_transform(dataset[feat])
          encoders[feat] = enc
      self.Encoders = encoders

      self.save_encoders(encoders)
      self.save_scaler(self.Scaler)
    
    return dataset

  def preprocess_to_predict_sli(self, original_dataset, target_sli):
    if len(original_dataset.shape) == 1:
      dataset = original_dataset.drop(target_sli)
    else:
      dataset = original_dataset.drop(target_sli, 1)
    return dataset

  def preprocess_to_predict_sli_to_knob(self, original_dataset, sli):

    knob_features = [k for k in original_dataset.keys() if "knob" in k]
    
    if len(original_dataset.shape) == 1:
      dataset = original_dataset.drop([r for r in original_dataset.keys() if (r in knob_features) or (r.startswith("app") and r != sli)])
    else:
      dataset = original_dataset.drop([r for r in original_dataset.keys() if (r in knob_features) or (r.startswith("app") and r != sli)], 1)

    return dataset


  def train_sli_to_knob_models(self):
    """
    Train models that will predict, given an sli, what is the best knob,
    for each of the considered knobs.
    """
    sli_knobs = {}
    
    non_knob_features = [k for k in self.Dataset.keys() if "knob" not in k]
    knob_features = [k for k in self.Dataset.keys() if "knob" in k]

    for sli in self.Dataset.keys():
      if "0.99" in sli:
        sli_knobs[sli] = {}
        for knob in knob_features:
          y = self.Dataset[knob]
          X = self.Dataset.drop([r for r in self.Dataset.keys() if (r in knob_features) or (r.startswith("app") and r != sli)], 1)
          
          regr = linear_model.LogisticRegression(C=.5, penalty='l1', tol=0.01, n_jobs=4)

          cv = ShuffleSplit(n_splits=3, test_size=0.2, random_state=42)
          
          scores = cross_val_score(regr, X, y, cv=cv)

          print("Accuracy for %s -> %s: %0.2f (+/- %0.2f)" % (sli, knob, scores.mean(), scores.std() * 2))
          
          # Now train the actual model with the whole dataset

          model = linear_model.LogisticRegression(C=.5, penalty='l1', tol=0.01, n_jobs=4) 

          model.fit(X, y)
          sli_knobs[sli][knob] = model

    return sli_knobs

  def train_sli_models(self):
    """
    Train sli models using the dataset being passed.
    Here the target is a SLI, we are trying to predict the SLI based on the other features.
    TODO: instead of `if 0.99 ...` we have to check if it is an actual sli.
    We could change the dataset to 'sli_houses_post'
    """
    sli_models = {}

    for target_feat in self.Dataset.keys():
      if "0.99" in target_feat:
        y = self.Dataset[target_feat]
        X = self.Dataset.drop(target_feat, 1)

        # Test accuracy with cross validation
        regr = linear_model.LassoCV()
        
        cv = ShuffleSplit(n_splits=3, test_size=0.2, random_state=42)
        scores = cross_val_score(regr, X, y, cv=cv)

        print("Accuracy for %s: %0.2f (+/- %0.2f)" % (target_feat, scores.mean(), scores.std() * 2))

        # Now train the actual model with the whole dataset
        model = linear_model.LassoCV() 
        model.fit(X, y)
        sli_models[target_feat] = model
    return sli_models
