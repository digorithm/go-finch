import pandas as pd
import sklearn
import matplotlib.pyplot as plt
import numpy as np
from sklearn import preprocessing
from pandas.plotting import scatter_matrix
from sklearn import linear_model
from sklearn.metrics import mean_squared_error, r2_score
from sklearn.model_selection import train_test_split
from sklearn.model_selection import learning_curve
from sklearn.model_selection import ShuffleSplit


def dict_compare(d1, d2):
  d1_keys = set(d1.keys())
  d2_keys = set(d2.keys())
  intersect_keys = d1_keys.intersection(d2_keys)
  added = d1_keys - d2_keys
  removed = d2_keys - d1_keys
  modified = {o : (d1[o], d2[o]) for o in intersect_keys if d1[o] != d2[o]}
  same = set(o for o in intersect_keys if d1[o] == d2[o])
  return added, removed, modified, same

def get_score(sli, predicted_knob_values, desired_configuration):
  if predicted_knob_values == desired_configuration:
    return 100.0
  else:
    _, _, diff, same = dict_compare(predicted_knob_values, desired_configuration)
    for wrong in diff:
      print("Failed to predict %s for SLI %s \n" % (wrong, sli))
      print("Expected %.2f, predicted %.2f \n" % (diff[wrong][1], diff[wrong][0]))
    return 100.0 - (len(diff)/len(desired_configuration) * 100)


def plot_learning_curve(estimator, title, X, y, variance=None, mse=None, ylim=None, cv=None,
                        n_jobs=1, train_sizes=np.linspace(.1, 1.0, 5)):
    """
    Generate a simple plot of the test and training learning curve.

    Parameters
    ----------
    estimator : object type that implements the "fit" and "predict" methods
        An object of that type which is cloned for each validation.

    title : string
        Title for the chart.

    X : array-like, shape (n_samples, n_features)
        Training vector, where n_samples is the number of samples and
        n_features is the number of features.

    y : array-like, shape (n_samples) or (n_samples, n_features), optional
        Target relative to X for classification or regression;
        None for unsupervised learning.

    ylim : tuple, shape (ymin, ymax), optional
        Defines minimum and maximum yvalues plotted.

    cv : int, cross-validation generator or an iterable, optional
        Determines the cross-validation splitting strategy.
        Possible inputs for cv are:
          - None, to use the default 3-fold cross-validation,
          - integer, to specify the number of folds.
          - An object to be used as a cross-validation generator.
          - An iterable yielding train/test splits.

        For integer/None inputs, if ``y`` is binary or multiclass,
        :class:`StratifiedKFold` used. If the estimator is not a classifier
        or if ``y`` is neither binary nor multiclass, :class:`KFold` is used.

        Refer :ref:`User Guide <cross_validation>` for the various
        cross-validators that can be used here.

    n_jobs : integer, optional
        Number of jobs to run in parallel (default 1).

    variance: float, coming from previous experiment
    mse: float, mean squared error, coming from previous experiment
    """
    plt.figure()
    plt.title(title)
    if ylim is not None:
        plt.ylim(*ylim)
    plt.xlabel("Training examples")
    plt.ylabel("Score")
    train_sizes, train_scores, test_scores = learning_curve(
        estimator, X, y, cv=cv, n_jobs=n_jobs, train_sizes=train_sizes)
    train_scores_mean = np.mean(train_scores, axis=1)
    train_scores_std = np.std(train_scores, axis=1)
    test_scores_mean = np.mean(test_scores, axis=1)
    test_scores_std = np.std(test_scores, axis=1)
    plt.grid()

    plt.fill_between(train_sizes, train_scores_mean - train_scores_std,
                     train_scores_mean + train_scores_std, alpha=0.1,
                     color="r")
    plt.fill_between(train_sizes, test_scores_mean - test_scores_std,
                     test_scores_mean + test_scores_std, alpha=0.1, color="g")
    plt.plot(train_sizes, train_scores_mean, 'o-', color="r",
             label="Training score")
    plt.plot(train_sizes, test_scores_mean, 'o-', color="g",
             label="Cross-validation score")

    if mse and variance:
      variance_label = "Variance: %.2f" % variance
      mse_label = "MSE: %.2f" % mse
      plt.plot([], [], ' ', label=[variance_label, mse_label])

    plt.legend(loc="best")
    plt.savefig(title + ".png")
    return plt

