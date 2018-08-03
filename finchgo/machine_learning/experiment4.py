import json
from finch_ml_test import validate_prediction_with_SLI_approach

# Loading SLO file" 
slos = json.load(open("slo.json"))

validate_prediction_with_SLI_approach(SLOs=slos, dataset_file="dataset_test.csv")