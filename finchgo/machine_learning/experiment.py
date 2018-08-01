import json
from finch_ml_test import validate_prediction_workflow

# Loading SLO file" 
slos = json.load(open("slo.json"))

# This configuration is based on what we know about optimal configuration
# If the predicted values are close to this, the prediction is good.
desired_configuration = {'pg_checkpoint_completion_target': 0.9,      
                        'pg_work_mem_kb': 3000,
                        'pg_maintenance_work_mem_mb': 1024, 'pg_shared_buffers_mb': 4000, 'pg_random_page_cost': 4, 'pg_wal_buffers_mb': 128, 'pg_default_statistics_target': 100, 'pg_effective_cache_size_mb': 120000,
                        'pg_max_wal_size_gb': 2,
                        'pg_min_wal_size_mb': 1000}


validate_prediction_workflow(slos=slos, desired_configuration=desired_configuration, dataset_file="dataset_test.csv")