import pandas as pd
import sklearn
import matplotlib.pyplot as plt
from sklearn import preprocessing
from pandas.plotting import scatter_matrix
from sklearn import linear_model
from sklearn.metrics import mean_squared_error, r2_score
from sklearn.model_selection import train_test_split
from sklearn.model_selection import learning_curve
from sklearn.model_selection import ShuffleSplit
from util import plot_learning_curve



def run_experiment():

  data = pd.read_csv("dataset.csv")

  # remove weird column
  data = data.drop("Unnamed: 38", 1)

  # remove timestamp, for now
  data = data.drop("timestamp", 1)

  # Encode the labels for pg features
  enc = preprocessing.LabelEncoder()

  for feat in data.keys():
      if feat.startswith("pg"):
          data[feat] = enc.fit_transform(data[feat])

  metrics = ["io_wait", "memory_usage", "disk_write_bytes", "cpu_usage", "disk_read_bytes", "cpu_idle"]

  # Train a model for each system's metric and report performance
  for metric in metrics:
      y = data[metric]
      X = data.drop(metric, 1)

      # Finding the MSE and R^2
      X_train, X_test, y_train, y_test = train_test_split(X, y, test_size=0.33, random_state=42)

      regr = linear_model.LinearRegression()
      regr.fit(X_train, y_train)
      y_pred = regr.predict(X_test)
      print("### Model for: %s ###" % metric)
      print("Mean squared error: %.2f" % mean_squared_error(y_test, y_pred))
      print("Variance: %.2f \n\n" % r2_score(y_test, y_pred))

      # Plotting learning curves through cross-validation
      title = "Learning curves for " + metric
      # Cross validation with 100 iterations to get smoother mean test and train
      # score curves, each time with 20% data randomly selected as a validation set.
      cv = ShuffleSplit(n_splits=100, test_size=0.2, random_state=42)

      estimator = linear_model.LinearRegression()

      p = plot_learning_curve(estimator, title, X, y, r2_score(y_test, y_pred), mean_squared_error(y_test, y_pred), ylim=(0.7, 1.01), cv=cv, n_jobs=4)


run_experiment()
