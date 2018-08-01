import React, { Component } from 'react';
import { api } from '../utils';
import './loginForm.css';

export default class LoginForm extends Component {

  constructor(props) {
    super(props);
    this.state = {
      username: '',
      password: ''
    }
  }

  usernameChange = event => {
    this.setState({username: event.target.value});
  }

  passwordChange = event => {
    this.setState({password: event.target.value});
  }

  handleSubmit = event => {
    const { username, password } = this.state;
    api.login(username, password)
      .then(res => this.props.setCredentials(username, password))
      .catch(err => {
        console.log(err);
        this.props.setError('Username and password did not match');
      });
    event.preventDefault();
    return false;
  }

  render() {
    const { username, password } = this.state;
    return (
      <div className='login-form'>
        <form onSubmit={this.handleSubmit}>
          <input
            type="text"
            value={username}
            placeholder="Username"
            onChange={this.usernameChange}
            className='card form-input form-component'
            autoFocus
          />
          <input
            type="password"
            value={password}
            placeholder="Password"
            onChange={this.passwordChange}
            className='card form-input form-component'
          />
        <input
          type="submit"
          value="Log In"
          className='card form-button form-component'/>
        </form>
      </div>
    )
  }

}
