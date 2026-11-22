## Chapter 2: Basic Operations ##

### Vectors, Data, Matrices and Subsetting ###

---

**In [1]:**

```
from __future__ import division
import pandas as pd
import numpy as np
import matplotlib.pyplot as plt
%matplotlib inline
```

---

**In [2]:**

```
x = np.array([2,7,5])    # explicit vector creation
x
```

**Out[2]:**

```
"array([2, 7, 5])"
```

---

**In [3]:**

```
y = np.arange(4, 13, 3)  # vector creation from a sequence (start, stop, step)
y
```

**Out[3]:**

```
"array([ 4,  7, 10])"
```

---

**In [4]:**

```
x + y    # vectors can be added
```

**Out[4]:**

```
"array([ 6, 14, 15])"
```

---

**In [5]:**

```
x / y    # divided
```

**Out[5]:**

```
"array([ 0.5,  1. ,  0.5])"
```

---

**In [6]:**

```
x ** y    # exponentiated
```

**Out[6]:**

```
"array([     16,  823543, 9765625])"
```

---

**In [7]:**

```
x[1]    # vector elements can be selected by position
```

**Out[7]:**

```
"7"
```

---

**In [8]:**

```
x[1:3]  # multiple elements can be selected using slices
```

**Out[8]:**

```
"array([7, 5])"
```

---

**In [9]:**

```
x[-2]  # elements can be specified as offset from end
```

**Out[9]:**

```
"7"
```

---

**In [10]:**

```
x[np.array([0,1])]  # elements can be specified as an array
```

**Out[10]:**

```
"array([2, 7])"
```

---

**In [11]:**

```
Z = np.matrix(np.arange(1,13)).reshape((4, 3))
Z                 # note: R arranges the elements column-wise
```

**Out[11]:**

```
"matrix([[ 1,  2,  3],\n        [ 4,  5,  6],\n        [ 7,  8,  9],\n        [10, 11, 12]])"
```

---

**In [12]:**

```
Z[2:4, 1:3]    # R is 1-based and includes ending index, Python is 0 based and does not.
```

**Out[12]:**

```
"matrix([[ 8,  9],\n        [11, 12]])"
```

---

**In [13]:**

```
Z[:, 1:3]    # column slice
```

**Out[13]:**

```
"matrix([[ 2,  3],\n        [ 5,  6],\n        [ 8,  9],\n        [11, 12]])"
```

---

**In [14]:**

```
Z.shape
```

**Out[14]:**

```
"(4, 3)"
```

---

### Generating Random Data, Graphics ###

---

**In [15]:**

```
x = np.random.uniform(0.0, 1.0, 50)
x
```

**Out[15]:**

```
"array([ 0.65268288,  0.55438063,  0.12332376,  0.73656161,  0.172588  ,\n        0.6952661 ,  0.36758793,  0.42034067,  0.54835434,  0.85493252,\n        0.36666796,  0.14029571,  0.16910043,  0.18956811,  0.77016189,\n        0.1914512 ,  0.48202498,  0.98936099,  0.43339457,  0.3690124 ,\n        0.91976653,  0.42354589,  0.23428098,  0.24846295,  0.20409963,\n        0.71024602,  0.57361799,  0.12510961,  0.70550836,  0.62120022,\n        0.7089459 ,  0.67670746,  0.39836928,  0.17109468,  0.57367169,\n        0.60584339,  0.72445482,  0.67991423,  0.94658844,  0.91068048,\n        0.65489981,  0.05985463,  0.56797862,  0.67604148,  0.989085  ,\n        0.79619792,  0.46406902,  0.91227664,  0.38043679,  0.04246386])"
```

---

**In [16]:**

```
y = np.random.normal(0.0, 1.0, 50)
y
```

**Out[16]:**

```
"array([ -5.95648785e-01,   1.54445571e+00,   1.35312925e+00,\n         1.88123098e+00,  -3.80353324e-01,  -9.67551972e-04,\n         5.97422145e-01,   3.76114342e-02,  -1.33779968e+00,\n        -5.95407697e-01,   1.27206236e+00,  -1.70114410e+00,\n        -1.84147571e+00,   7.16611666e-01,   2.33006825e+00,\n        -6.56709283e-01,   8.26951067e-01,   4.64854930e-01,\n         8.79290232e-01,  -7.42861574e-02,   4.88797828e-01,\n         1.50860241e+00,   1.07467313e+00,  -7.99843607e-01,\n         6.81719603e-02,   5.06424430e-01,   2.15770333e-01,\n         1.67843398e+00,  -5.01086886e-01,   1.51407152e-01,\n        -3.30648417e-02,  -8.50683506e-01,  -1.56531920e+00,\n        -6.43438018e-01,  -9.55514605e-01,   1.51696803e-01,\n        -9.87417666e-01,  -4.44820480e-01,   1.21902188e-01,\n        -5.20474443e-01,   8.31544982e-01,  -6.75527626e-01,\n         1.00537660e+00,   4.87948062e-01,   3.43123002e-01,\n         3.95408496e-01,   4.27046930e-01,   3.69945068e-01,\n         1.36728064e+00,  -3.15495180e-01])"
```

---

**In [17]:**

```
fig, ax = plt.subplots()
plt.scatter(x, y)
```

**Out[17]:**

```
"\u003cmatplotlib.collections.PathCollection at 0x420b410\u003e"
```

---

**In [18]:**

```
fig, ax = plt.subplots()
plt.xlabel("Random Uniform")
plt.ylabel("Random Normal")
plt.scatter(x, y, marker='o', color='red')  # plot customizations
```

**Out[18]:**

```
"\u003cmatplotlib.collections.PathCollection at 0x4241f50\u003e"
```

---

**In [19]:**

```
plt.subplot(121)    # parameter indicates 1 rows, 2 col, first figure
plt.scatter(x, y)
plt.subplot(122)
plt.hist(y)
```

**Out[19]:**

```
"(array([ 3.,  1.,  9.,  5.,  9.,  9.,  6.,  3.,  4.,  1.]),\n array([-1.84147571, -1.42432131, -1.00716692, -0.59001252, -0.17285812,\n        0.24429627,  0.66145067,  1.07860506,  1.49575946,  1.91291385,\n        2.33006825]),\n \u003ca list of 10 Patch objects\u003e)"
```

---

### Reading in data ###

---

**In [20]:**

```
# data comes from ISL book site: http://www-bcf.usc.edu/~gareth/ISL/data.html
auto_df = pd.read_csv("../data/Auto.csv")
auto_df.columns     # column names
```

**Out[20]:**

```
"Index([u'mpg', u'cylinders', u'displacement', u'horsepower', u'weight', u'acceleration', u'year', u'origin', u'name'], dtype='object')"
```

---

**In [21]:**

```
auto_df.shape    # number of rows, number of columns
```

**Out[21]:**

```
"(397, 9)"
```

---

**In [22]:**

```
type(auto_df)
```

**Out[22]:**

```
"pandas.core.frame.DataFrame"
```

---

**In [23]:**

```
auto_df.describe()  # equivalent of R's DataFrame.summary()
```

**Out[23]:**

```html
"\u003cdiv style=\"max-height:1000px;max-width:1500px;overflow:auto;\"\u003e\n\u003ctable border=\"1\" class=\"dataframe\"\u003e\n  \u003cthead\u003e\n    \u003ctr style=\"text-align: right;\"\u003e\n      \u003cth\u003e\u003c/th\u003e\n      \u003cth\u003empg\u003c/th\u003e\n      \u003cth\u003ecylinders\u003c/th\u003e\n      \u003cth\u003edisplacement\u003c/th\u003e\n      \u003cth\u003eweight\u003c/th\u003e\n      \u003cth\u003eacceleration\u003c/th\u003e\n      \u003cth\u003eyear\u003c/th\u003e\n      \u003cth\u003eorigin\u003c/th\u003e\n    \u003c/tr\u003e\n  \u003c/thead\u003e\n  \u003ctbody\u003e\n    \u003ctr\u003e\n      \u003cth\u003ecount\u003c/th\u003e\n      \u003ctd\u003e 397.000000\u003c/td\u003e\n      \u003ctd\u003e 397.000000\u003c/td\u003e\n      \u003ctd\u003e 397.000000\u003c/td\u003e\n      \u003ctd\u003e  397.000000\u003c/td\u003e\n      \u003ctd\u003e 397.000000\u003c/td\u003e\n      \u003ctd\u003e 397.000000\u003c/td\u003e\n      \u003ctd\u003e 397.000000\u003c/td\u003e\n    \u003c/tr\u003e\n    \u003ctr\u003e\n      \u003cth\u003emean\u003c/th\u003e\n      \u003ctd\u003e  23.515869\u003c/td\u003e\n      \u003ctd\u003e   5.458438\u003c/td\u003e\n      \u003ctd\u003e 193.532746\u003c/td\u003e\n      \u003ctd\u003e 2970.261965\u003c/td\u003e\n      \u003ctd\u003e  15.555668\u003c/td\u003e\n      \u003ctd\u003e  75.994962\u003c/td\u003e\n      \u003ctd\u003e   1.574307\u003c/td\u003e\n    \u003c/tr\u003e\n    \u003ctr\u003e\n      \u003cth\u003estd\u003c/th\u003e\n      \u003ctd\u003e   7.825804\u003c/td\u003e\n      \u003ctd\u003e   1.701577\u003c/td\u003e\n      \u003ctd\u003e 104.379583\u003c/td\u003e\n      \u003ctd\u003e  847.904119\u003c/td\u003e\n      \u003ctd\u003e   2.749995\u003c/td\u003e\n      \u003ctd\u003e   3.690005\u003c/td\u003e\n      \u003ctd\u003e   0.802549\u003c/td\u003e\n    \u003c/tr\u003e\n    \u003ctr\u003e\n      \u003cth\u003emin\u003c/th\u003e\n      \u003ctd\u003e   9.000000\u003c/td\u003e\n      \u003ctd\u003e   3.000000\u003c/td\u003e\n      \u003ctd\u003e  68.000000\u003c/td\u003e\n      \u003ctd\u003e 1613.000000\u003c/td\u003e\n      \u003ctd\u003e   8.000000\u003c/td\u003e\n      \u003ctd\u003e  70.000000\u003c/td\u003e\n      \u003ctd\u003e   1.000000\u003c/td\u003e\n    \u003c/tr\u003e\n    \u003ctr\u003e\n      \u003cth\u003e25%\u003c/th\u003e\n      \u003ctd\u003e  17.500000\u003c/td\u003e\n      \u003ctd\u003e   4.000000\u003c/td\u003e\n      \u003ctd\u003e 104.000000\u003c/td\u003e\n      \u003ctd\u003e 2223.000000\u003c/td\u003e\n      \u003ctd\u003e  13.800000\u003c/td\u003e\n      \u003ctd\u003e  73.000000\u003c/td\u003e\n      \u003ctd\u003e   1.000000\u003c/td\u003e\n    \u003c/tr\u003e\n    \u003ctr\u003e\n      \u003cth\u003e50%\u003c/th\u003e\n      \u003ctd\u003e  23.000000\u003c/td\u003e\n      \u003ctd\u003e   4.000000\u003c/td\u003e\n      \u003ctd\u003e 146.000000\u003c/td\u003e\n      \u003ctd\u003e 2800.000000\u003c/td\u003e\n      \u003ctd\u003e  15.500000\u003c/td\u003e\n      \u003ctd\u003e  76.000000\u003c/td\u003e\n      \u003ctd\u003e   1.000000\u003c/td\u003e\n    \u003c/tr\u003e\n    \u003ctr\u003e\n      \u003cth\u003e75%\u003c/th\u003e\n      \u003ctd\u003e  29.000000\u003c/td\u003e\n      \u003ctd\u003e   8.000000\u003c/td\u003e\n      \u003ctd\u003e 262.000000\u003c/td\u003e\n      \u003ctd\u003e 3609.000000\u003c/td\u003e\n      \u003ctd\u003e  17.100000\u003c/td\u003e\n      \u003ctd\u003e  79.000000\u003c/td\u003e\n      \u003ctd\u003e   2.000000\u003c/td\u003e\n    \u003c/tr\u003e\n    \u003ctr\u003e\n      \u003cth\u003emax\u003c/th\u003e\n      \u003ctd\u003e  46.600000\u003c/td\u003e\n      \u003ctd\u003e   8.000000\u003c/td\u003e\n      \u003ctd\u003e 455.000000\u003c/td\u003e\n      \u003ctd\u003e 5140.000000\u003c/td\u003e\n      \u003ctd\u003e  24.800000\u003c/td\u003e\n      \u003ctd\u003e  82.000000\u003c/td\u003e\n      \u003ctd\u003e   3.000000\u003c/td\u003e\n    \u003c/tr\u003e\n  \u003c/tbody\u003e\n\u003c/table\u003e\n\u003cp\u003e8 rows × 7 columns\u003c/p\u003e\n\u003c/div\u003e"
```

---

**In [24]:**

```
plt.ylabel("MPG")
auto_df.plot(x="cylinders", y="mpg", style='o')
```

**Out[24]:**

```
"\u003cmatplotlib.axes.AxesSubplot at 0x46ef650\u003e"
```

---

**In [25]:**

```
auto_df.boxplot(column="mpg", by="cylinders") # MPG distribution by number of cylinders
```

**Out[25]:**

```
"\u003cmatplotlib.axes.AxesSubplot at 0x499f510\u003e"
```

---

**In [26]:**

```
# similar to R pairs, shows correlation scatter plots between columns and distribution for each 
# column along the diagonal.
# The R version that uses formulas does not seem to have a Python equivalent (and doesn't seem
# to be very useful for exploratory analysis IMO).
axes = pd.tools.plotting.scatter_matrix(auto_df, color="brown")
```
