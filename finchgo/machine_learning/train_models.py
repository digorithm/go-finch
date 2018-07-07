"""
Simple script to call the method that will train all models. The dataset used will be the one being copied from FinchGO to here. It's called 'dataset.csv'
"""

from finch_ml import FinchML

finch = FinchML()

finch.train_models()
