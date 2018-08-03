import json
from finch_ml_test import validate_prediction_for_single_sla_workflow

# Loading SLO file" 
slos = json.load(open("slo.json"))

# We want to predict the optimal knobs to optimize THIS single SLA
violatedSLA = 'app_http_request_latency_POST_users_0.99'
violatedSLO = {}
violatedSLO[violatedSLA] = slos[violatedSLA]

validate_prediction_for_single_sla_workflow(violatedSLO=violatedSLO, dataset_file="dataset_test.csv")