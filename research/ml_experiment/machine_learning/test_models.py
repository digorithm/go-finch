"""
Module responsible for testing the performance of models for a given dataset
"""

import pandas as pd
from models import train_models, preprocess_to_predict_sli_to_knob, preprocess_to_predict_sli
from collections import Counter

def dict_compare(d1, d2):
  d1_keys = set(d1.keys())
  d2_keys = set(d2.keys())
  intersect_keys = d1_keys.intersection(d2_keys)
  added = d1_keys - d2_keys
  removed = d2_keys - d1_keys
  modified = {o : (d1[o], d2[o]) for o in intersect_keys if d1[o] != d2[o]}
  same = set(o for o in intersect_keys if d1[o] == d2[o])
  return added, removed, modified, same

def get_score(sli, previous_knob_values,
              predicted_knob_values, desired_configuration):
  if predicted_knob_values == desired_configuration:
    return 100.0
  else:
    _, _, diff, same = dict_compare(predicted_knob_values, desired_configuration)
    for wrong in diff:
      print("Failed to predict %s for SLI %s \n" % (wrong, sli))
      print("Expected %.2f, predicted %.2f \n" % (diff[wrong][1], diff[wrong][0]))
    return 100.0 - (len(diff)/len(desired_configuration) * 100)

def validate_prediction_workflow(slos, desired_configuration, dataset_file="dataset.csv"):
  """
  Validate the prediction models by performing bidirectional predictions
  slos is a dictionary where keys are the SLIs and the values are its respective SLOs
  """
  dataset = pd.read_csv(dataset_file)

  ds = dataset.loc[1:215]

  # This datapoint is a sample of a bad performance. All SLIs are violating our fictional SLOs
  test = dataset.loc[0]

  # Train the model using the dataset that does not contain the test sample
  sli_knob_models, sli_models, encoders = train_models(dataframe=ds)

  # Grab knob names
  knobs = [pg for pg in dataset.keys() if pg.startswith("pg")]

  all_predictions = {}
  for sli, slo in slos.items():
    # Grab the values in the knobs of our test datapoint
    previous_knob_values = {}
    for knob in knobs:
      previous_knob_values[knob] = test[knob]

    # Let's test with this single SLI. Our goal is to predict all the opposite knobs
    sli_test = preprocess_to_predict_sli_to_knob(test, sli=sli)

    # Our desired SLI value. We change the test sample to contain the desired SLI but keeping the current system's context
    sli_test[sli] = slo

    predicted_knob_values = {}
    for knob in knobs:
      pred = sli_knob_models[sli][knob].predict([sli_test.as_matrix()])
      predicted_knob_values[knob] = encoders[knob].inverse_transform(pred)[0]

    all_predictions[sli] = predicted_knob_values

    score = get_score(sli, previous_knob_values,predicted_knob_values, desired_configuration)
    print("Score for sli %s: %.2f" % (sli, score)) 


  final_knob_predictions = {}
  for knob in knobs:
    # Voting mechanism, it grabs the highest occurence among all predictions for this given knob
    final_knob_predictions[knob] = Counter([all_predictions[sli][knob] for sli in slos.keys()]).most_common(1)[0][0]

  # Now we go the other way around: we predict the SLI using the set of predicted knobs, just to validate the previous prediction
  inverse_test = test
  for knob in final_knob_predictions:
    inverse_test[knob] = final_knob_predictions[knob]

  # Predicing the new SLI given the predicted knobs
  for sli in slos.keys():
    inverse_test_without_sli = preprocess_to_predict_sli(inverse_test, sli)
    slis_to_be_modified = [s for s in slos.keys() if s != sli]
    # Pulling the other SLIs to be equal to SLO
    for s in slis_to_be_modified:
      inverse_test_without_sli[s] = slos[s]

    pred = sli_models[sli].predict([inverse_test_without_sli])
    import pdb; pdb.set_trace()
    print("### \n SLI: %s \n SLO: %.f \n Predicted: %.f \n difference between prediction and SLO: %.f \n### \n" % (sli, slos[sli], pred[0], slos[sli] - pred[0]))
