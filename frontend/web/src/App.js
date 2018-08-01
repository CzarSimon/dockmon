import React, { Component } from 'react';
import AppHeading from './components/appHeading';
import ServiceList from './components/serviceList';
import LoginForm from './components/loginForm';
import ErrorMessage from './components/errorMessage';
import { orientation } from './utils';
import './App.css';

export default class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      orientationStyle: orientation.getStyle(),
    }

    window.addEventListener('resize', event => {
      this.setState({
        orientationStyle: orientation.getStyle(),
      })
    })
  }

  setCredentials = (username, password) => {
    this.setState({
      error: null,
      username,
      password,
    });
  }

  unsetCredentials = () => {
    this.setState({
      error: null,
      username: null,
      password: null
    });
  }

  setError = errorMessage => {
    this.setState({
      error: errorMessage
    });
  }

  loggedIn = () => (this.state.username && this.state.password)

  render() {
    const { username, password, error } = this.state;
    return (
      <div className='app'>
        <AppHeading />
        <div className='content' style={this.state.orientationStyle}>
          <ErrorMessage text={error} />
          {
            (this.loggedIn())
              ? <ServiceList username={username} password={password} unsetCredentials={this.unsetCredentials}/>
              : <LoginForm setCredentials={this.setCredentials} setError={this.setError} />
          }
        </div>
      </div>
    )
  }
}
