import React, { Component } from 'react';

export default class LogoutButton extends Component {
  handleClick = event => {
    this.props.unsetCredentials();
    event.preventDefault();
    return false;
  }

  render() {
    return (
      <button onClick={this.handleClick} className='card form-button form-component'>
        Log out
      </button>
    );
  }
}
