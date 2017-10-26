from japronto import Application, RouteNotFoundException
from scipy import sparse
import json
import numpy as np
import redis
import pickle

r = redis.StrictRedis(host='localhost', port=6379, db=0)
sm_features_column_length = int(r.get("smart_clf_features_column_length"))

try:
    smart_clf = pickle.loads(r.get("smart_clf_browser"))
except:
    print("Cant load model from smart_clf redis storage")

cdef predict(request):
    features_list = request.query["positions"].split(',')
    features_list = [int(a) for a in features_list if a not in ""]

    rows = []
    cols = []
    data = []

    for index in range(sm_features_column_length):
        rows.append(0)
        cols.append(index)
        if index in features_list:
            data.append(1)
        else:
            data.append(0)

    x_test = sparse.csr_matrix((data, (rows, cols)), dtype=np.int8)
    predict_proba = smart_clf.predict_proba(x_test)
    results = []
    i = 0

    for idx, val in enumerate(smart_clf.classes_):
        _new = {}
        _new[val] = round(predict_proba[0][idx], 10)
        results.append(_new)

    return request.Response(text=json.dumps(results))

cdef reload_model(request):
    smart_clf_features_column_length = int(r.get("smart_clf_features_column_length"))
    smart_clf = pickle.loads(r.get("smart_clf_browser"))
    print("Model was reloaded")
    print("smart_clf_features_column_length is", smart_clf_features_column_length)
    return request.Response(text="reloaded")

app = Application()
app.router.add_route('/', predict)
app.router.add_route('/reload', reload_model)
app.run(debug=False, port=8081)
