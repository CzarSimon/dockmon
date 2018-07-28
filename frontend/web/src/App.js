import React, { Component } from 'react';
import AppHeading from './components/appHeading';
import ServiceList from './components/serviceList';
import { orientation } from './utils';
import './App.css';

export default class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      orientationStyle: orientation.getStyle()
    }

    window.addEventListener('resize', event => {
      this.setState({
        orientationStyle: orientation.getStyle(),
      })
    })
  }

  render() {
    return (
      <div className='app'>
        <AppHeading />
        <div className='content' style={this.state.orientationStyle}>
          <ServiceList />
        </div>
      </div>
    )
  }
}
