"""
Module responsible for testing the performance of models for a given dataset
"""

import pandas as pd
from finch_ml import FinchML
from collections import Counter
from util import get_score


def validate_prediction_workflow(slos, desired_configuration, dataset_file="dataset.csv"):
  """
  Validate the prediction models by performing bidirectional predictions
  slos is a dictionary where keys are the SLIs and the values are its respective SLOs
  """
  dataset = pd.read_csv(dataset_file)
  #dataset = dataset.drop(["io_wait", "memory_usage", "cpu_usage", "cpu_idle", "disk_write_bytes", "disk_read_bytes"], 1)
  
  ds = dataset.loc[1:].copy()

  # Train the model using the dataset that does not contain the test sample
  finch = FinchML()
  finch.train_models(dataframe=ds)

  # This datapoint is a sample of a bad performance. All SLIs are violating our fictional SLOs
  test = dataset.loc[0].copy()

  # Grab knob names
  knobs = [pg for pg in dataset.keys() if pg.startswith("pg")]
  final_knob_predictions = finch.predict_optimal_knobs(test, knobs, slos)

  sli_predictions = finch.predict_SLIs(test, slos, predicted_knobs=final_knob_predictions)
  
  original_scale = dataset.loc[0].copy()
  original_scale = finch.preprocess_data(original_scale)
  for sli in sli_predictions.keys():
    original_scale[sli] = sli_predictions[sli]
  
  # Scale back to see results
  original_scale[[k for k in original_scale.keys() if not k.startswith("pg")]] = finch.Scaler.inverse_transform(original_scale[[k for k in original_scale.keys() if not k.startswith("pg")]])

  scaled_back_sli_predictions = {}

  for sli in sli_predictions.keys():
    scaled_back_sli_predictions[sli] = original_scale[sli]


  for sli in slos.keys():
    print("### \n SLI: %s \n SLO: %.f \n Predicted: %.f \n difference between prediction and SLO: %.f \n### \n" % (sli, slos[sli], scaled_back_sli_predictions[sli], slos[sli] - scaled_back_sli_predictions[sli]))

  print("Optimal set of knobs for test case:: " + str(final_knob_predictions))
