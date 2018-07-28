import React, { Component } from 'react';
import ServiceHeading from './serviceHeading';
import ServiceInfo from './serviceInfo';
import './serviceStatus.css';

export default class ServiceStatus extends Component {
  state = {
    clicked: false
  }

  handleClicked = event => {
    this.setState({
      clicked: !this.state.clicked
    })
  }

  render() {
    return (
      <div onClick={this.handleClicked} className="service-status card">
        <ServiceHeading {...this.props} />
        {(this.state.clicked) ? <ServiceInfo {...this.props} /> : <div/>}
      </div>
    )
  }

}
