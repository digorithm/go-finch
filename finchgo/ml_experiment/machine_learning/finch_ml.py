"""
This module will contain all models
"""

import pandas as pd
import sklearn
from sklearn import preprocessing
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

    return sli_knob_models, sli_models
  
  def predict_optimal_knobs(self, X, knobs, slos):
    """
    Given the system's context X, the used knobs, and the defined SLOs,
    predict the optimal set of knob values
    """

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

    # Unique sample
    if len(original_dataset.shape) == 1:
      # remove weird column
      dataset = original_dataset.drop("Unnamed: 39")

      # remove timestamp, for now
      dataset = dataset.drop("timestamp")
      scaled = self.Scaler.transform([dataset[[k for k in dataset.keys() if not k.startswith("pg")]]])[0]
      dataset[[k for k in dataset.keys() if not k.startswith("pg")]] = scaled
      
      for feat in dataset.keys():
        if feat.startswith("pg"):
          # Encode the labels for pg features
          dataset[feat] = self.Encoders[feat].transform([dataset[feat]])

    # Training dataset  
    else:
      # remove weird column
      dataset = original_dataset.drop("Unnamed: 39", 1)

      # remove timestamp, for now
      dataset = dataset.drop("timestamp", 1)
      
      # Scale only features that are not pg variables
      self.Scaler = preprocessing.StandardScaler().fit(dataset[[k for k in dataset.keys() if not k.startswith("pg")]])

      dataset[[k for k in dataset.keys() if not k.startswith("pg")]] = self.Scaler.transform(dataset[[k for k in dataset.keys() if not k.startswith("pg")]])
      # Use encoders dict whenever you need to retrieve the real predicted label
      encoders = {}
      for feat in dataset.keys():
        if feat.startswith("pg"):
          # Encode the labels for pg features
          enc = preprocessing.LabelEncoder()
          dataset[feat] = enc.fit_transform(dataset[feat])
          encoders[feat] = enc
      self.Encoders = encoders

    
    return dataset

  def preprocess_to_predict_sli(self, original_dataset, target_sli):
    if len(original_dataset.shape) == 1:
      dataset = original_dataset.drop(target_sli)
    else:
      dataset = original_dataset.drop(target_sli, 1)
    return dataset

  def preprocess_to_predict_sli_to_knob(self, original_dataset, sli):
    
    if len(original_dataset.shape) == 1:
      dataset = original_dataset.drop([r for r in original_dataset.keys() if r.startswith("pg") or (r.startswith("app") and r != sli)])
    else:
      dataset = original_dataset.drop([r for r in original_dataset.keys() if r.startswith("pg") or (r.startswith("app") and r != sli)], 1)

    return dataset


  def train_sli_to_knob_models(self):
    """
    Train models that will predict, given an sli, what is the best knob,
    for each of the considered knobs.
    """
    sli_knobs = {}

    for sli in self.Dataset.keys():
      if "0.99" in sli:
        sli_knobs[sli] = {}
        for knob in self.Dataset.keys():
          if knob.startswith("pg"):
        
            y = self.Dataset[knob]
            X = self.Dataset.drop([r for r in self.Dataset.keys() if r.startswith("pg") or (r.startswith("app") and r != sli)], 1)
            
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
