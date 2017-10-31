import time
import os
import sys
import numpy as np
cimport numpy as np
import redis
import pickle
from scipy import sparse
from sklearn.linear_model import LogisticRegression
from sklearn.multiclass import OneVsRestClassifier

import psycopg2
import warnings

warnings.filterwarnings("ignore")

os.environ["JOBLIB_TEMP_FOLDER"] = "/tmp"
#os.environ["JOBLIB_TEMP_FOLDER"] = "/media/levabd/ScienceProjects/data"
try:
    conn = psycopg2.connect("dbname='antifraud' user='antifraud' host='localhost' password=''")
except:
    print ("I am unable to connect to the database")


class SimpleProgressBar(object):
    def __init__(self, maximum, state=0):
        self.max = maximum
        self.state = state

    def _carriage_return(self):
        sys.stdout.write('\r')
        sys.stdout.flush()

    def _display(self):
        stars = ''.join(['*'] * self.state + [' '] * (self.max - self.state))
        print ('[{0}] {1}/{2}'.format(stars, self.state, self.max))
        self._carriage_return()

    def update(self, value=None):
        if not value is None:
            self.state = value
        self._display()

cpdef run_education():
    cdef int start = time.time()

    print("JOBLIB_TEMP_FOLDER: ", os.environ["JOBLIB_TEMP_FOLDER"])

    cursor = conn.cursor()
    cursor.execute("""select * from features; """)

    rows = []
    cols = []
    data = []
    for record in cursor.fetchall():
        rows.append(record[1])
        cols.append(record[2])
        data.append(1)

    sm_features= sparse.csr_matrix((data, (rows, cols)), dtype=np.int8)

    if sm_features == None:
        print("Cant educate - cant establish internet connection")
        return

    cursor = conn.cursor()
    cursor.execute("""select * from browsers;""")
    browsers =[]
    for record in cursor.fetchall():
        browsers.append(record[1])

    try:
        r = redis.StrictRedis(host='localhost', port=6379, db=0)
    except:
        print("Unable to establish redis connection")

    print("Education started")

    smart_clf = OneVsRestClassifier(LogisticRegression(C=100, n_jobs=-1), n_jobs=-1)
    smart_clf.fit(sm_features, browsers)

    r.set("smart_clf_features_column_length", sm_features.shape[1])
    r.set("smart_clf_browser", pickle.dumps(smart_clf))
    conn.close()

    print("Education finished")
    print("Model education took {} seconds".format(time.time() - start))

if __name__ == '__main__':
    run_education()
