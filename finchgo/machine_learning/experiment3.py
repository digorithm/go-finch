import json
from finch_ml_test import validate_prediction_with_multioutput_regressor_workflow

# Loading SLO file" 
slos = json.load(open("slo.json"))

validate_prediction_with_multioutput_regressor_workflow(SLOs=slos, dataset_file="dataset_test.csv")
