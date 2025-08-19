from flask import Flask, request, jsonify
import requests
import pandas as pd
from flask_cors import CORS
import joblib
import base64
import io
from sklearn.linear_model import LogisticRegression
from sklearn.tree import DecisionTreeClassifier
from sklearn.ensemble import RandomForestClassifier
from sklearn.svm import SVC
from sklearn.neighbors import KNeighborsClassifier
from sklearn.naive_bayes import GaussianNB
import xgboost as xgb
from sklearn.preprocessing import MinMaxScaler


# تعریف شیء app
app = Flask(__name__)
CORS(app)  # فعال‌سازی CORS برای تمامی درخواست‌ها

models = {
    'Linear Regression': LogisticRegression(penalty='l1', solver='liblinear', C=1.0, max_iter=10000, random_state=42),
    'Decision Tree': DecisionTreeClassifier(random_state=42),
    'Random Forest': RandomForestClassifier(n_estimators=100, random_state=42),
    'SVM': SVC(kernel='linear', random_state=42),
    'KNN': KNeighborsClassifier(),
    'Naive Bayes': GaussianNB(),
    'XGBoost': xgb.XGBClassifier(random_state=42)
}
@app.route('/predict', methods=['POST'])
def predict():
    try:
        data = request.get_json()
        # دریافت لیست دیتاست‌ها
        train_datasets = data['train_datasets']
        test_datasets = data['test_datasets']
        model_name = data['model']
        pre_built_model = data['is_custom_model']
        custom_model = data.get('custom_model')  # دریافت مدل سفارشی از Base64

        # استخراج شناسه‌ها از دیتاست‌ها
        train_dataset_ids = [dataset.split(' - ')[0] for dataset in train_datasets]
        test_dataset_ids = [dataset.split(' - ')[0] for dataset in test_datasets]

        # بارگذاری دیتاست‌ها از ID‌ها
        train_data = get_datasets_by_ids(train_dataset_ids)
        test_data = get_datasets_by_ids(test_dataset_ids)

        # بررسی وجود ستون 'bug' در دیتاست‌ها
        if 'bug' not in train_data.columns or 'bug' not in test_data.columns:
            return jsonify({'error': "'bug' column is missing in the dataset"}), 400

        if not pre_built_model:
            # اگر مدل پیش‌ساخته نیست، داده‌ها را ترکیب می‌کنیم
            X_train = train_data.select_dtypes(include=['float64', 'int64']).drop(
          columns=[col for col in ["version", "bug"] if col in train_data.select_dtypes(
      include=['float64', 'int64']).columns])  # فرض می‌کنیم که ستون 'bug' هدف است
            y_train = train_data['bug']

            # انتخاب مدل
            model = models[model_name]
            trained_model = model.fit(X_train, y_train)

            # پیش‌بینی
            X_test = test_data[X_train.columns]
            y_test = test_data['bug']
            preds = trained_model.predict(X_test)

        else:
            # اگر مدل پیش‌ساخته است، از آن استفاده می‌کنیم
            if custom_model:
                custom_model_bytes = base64.b64decode(custom_model)  # تبدیل Base64 به باینری
                custom_model = joblib.load(io.BytesIO(custom_model_bytes))  # بارگذاری مدل از بایت‌ها
                X_test = test_data.drop(columns=['bug'])
                preds = custom_model.predict(X_test)
            else:
                return jsonify({'error': 'Custom model is missing'}), 400

        # ایجاد پاسخ
        result = test_data.copy()
        result['predictions'] = preds

        # تبدیل پاسخ به فرمت مناسب برای ارسال به Flutter
        result_dict = result.to_dict(orient='records')
        return jsonify({'result': result_dict, 'accuracy': (preds == y_test).mean()}), 200

    except Exception as e:
        print(f"Error: {str(e)}")  # چاپ خطا برای رفع مشکل
        return jsonify({'error': str(e)}), 500

# این تابع برای بارگذاری دیتاست‌ها از ID‌ها است
def get_datasets_by_ids(dataset_ids):
    datasets = []
    for dataset_id in dataset_ids:
        url = f"http://localhost:8081/dataset/{dataset_id}"  # URL سرور Go برای دریافت دیتاست
        response = requests.get(url)
        
        if response.status_code == 200:
            dataset_content_base64 = response.json().get('content')  # دریافت داده‌های Base64 از پاسخ
            dataset_content = base64.b64decode(dataset_content_base64)  # تبدیل Base64 به داده‌های باینری
            datasets.append(pd.read_csv(io.BytesIO(dataset_content)))  # بارگذاری داده‌ها از بایت‌ها به صورت DataFrame
        else:
            print(f"Error fetching dataset {dataset_id}: {response.status_code}, {response.text}")
    
    if len(datasets) > 1 :
        return pd.concat(datasets, axis=0)
    else:
        return datasets[0]

if __name__ == '__main__':
    app.run(debug=True, host='0.0.0.0', port=5001)
