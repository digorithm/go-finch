from finch_ml import FinchML
import pandas as pd
import json
import sys

SLOs = json.load(open("slo.json"))

# Here we use loc[0] because it turns it into Series instead of
# a full dataframe
X = pd.read_csv("single.csv").loc[0]

# If storage/schedule endpoint haven't been called yet
# we fill it with expected values. This should be temp
late_keys = {"app_http_request_latency_POST_schedules_0.5": 50.0, "app_http_request_latency_POST_schedules_0.9": 700.0, "app_http_request_latency_POST_schedules_0.99": 1500.0, "app_http_request_latency_POST_storages_0.5": 1800.0, "app_http_request_latency_POST_storages_0.9": 3200.0, "app_http_request_latency_POST_storages_0.99": 4000.0}

for key, value in late_keys.items():
  if key not in X:
    X[key] = value

knobs = [k for k in X.keys() if "knob" in k]

finch = FinchML()

final_knob_predictions = finch.predict_optimal_knobs(X, knobs, SLOs)


# For now we are doing this: the knob prediction come as 
# "app_knob_<KNOB name>", but we want to return just the knob name
# to the Go component, so we split it by '_' and get the third element
final_knob_predictions_correct_form = {}
for knob, value in final_knob_predictions.items():
  correct_knob_name = knob.split('_')[2]
  final_knob_predictions_correct_form[correct_knob_name] = value

print(json.dumps(final_knob_predictions_correct_form))
sys.stdout.flush()