import React, { Component } from 'react';
import { api } from '../utils';
import ServiceStatus from './serviceStatus';

export default class ServiceList extends Component {

  constructor(props) {
    super(props);
    this.state = {
      services: null
    }
  }

  componentDidMount() {
    this.fetchServiceState();
    setInterval(this.fetchServiceState, 10 * 1000);
  }

  fetchServiceState = () => {
    api.getServiceStatuses('token')
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
      <div className='service-list'>
        {services.map((service, i) => <ServiceStatus key={i} {...service} />)}
      </div>
    )
  }

}
