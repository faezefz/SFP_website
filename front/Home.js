import React, { useEffect } from 'react';
import { useHistory } from 'react-router-dom';

const Home = () => {
  const history = useHistory();

  useEffect(() => {
    const token = localStorage.getItem('token'); // گرفتن توکن از localStorage

    if (token) {
      // اگر توکن موجود بود، به داشبورد هدایت می‌کنیم
      history.push(`/dashboard/1`); // اینجا باید شناسه کاربر رو از توکن استخراج کنی
    } else {
      // اگر توکن نباشه، به صفحه ورود هدایت می‌کنیم
      history.push('/login');
    }
  }, [history]);

  return (
    <div>
      <h1>Welcome to the Home page</h1>
    </div>
  );
};

export default Home;
