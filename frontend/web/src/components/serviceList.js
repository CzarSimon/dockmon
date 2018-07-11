import React, { Component } from 'react';
import { getServiceStatuses } from '../api';
import ServiceStatus from './serviceStatus';

export default class ServiceList extends Component {

  constructor(props) {
    super(props);
    this.state = {
      services: null
    }
  }

  componentDidMount() {
    getServiceStatuses('token')
      .then(services => this.setState({ services: services }))
  }

  render() {
    const { services } = this.state;
    if (!services) {
      return (
        <div />
      )
    }
    return (
      <div>
        {services.map((service, i) => <ServiceStatus key={i} {...service} />)}
      </div>
    )
  }

}
