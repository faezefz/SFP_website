import React from 'react';
import { BrowserRouter as Router, Route, Switch } from 'react-router-dom';
import Home from './Home'; // صفحه خانه
import Login from './Login'; // صفحه ورود
import Dashboard from './Dashboard'; // صفحه داشبورد

const App = () => {
  return (
    <Router>
      <Switch>
        <Route exact path="/" component={Home} />
        <Route path="/login" component={Login} />
        <Route path="/dashboard/:id" component={Dashboard} />
      </Switch>
    </Router>
  );
};

export default App;
