import React, { Component } from 'react';
import { api } from '../utils';
import ServiceStatus from './serviceStatus';
import LogoutButton from './logoutButton';

export default class ServiceList extends Component {

  constructor(props) {
    super(props);
    this.state = {
      services: null,
      intervalId: null
    }
  }

  componentDidMount() {
    this.fetchServiceState();
    const intervalId = setInterval(this.fetchServiceState, 10 * 1000);
    this.setState({ intervalId });
  }

  fetchServiceState = () => {
    const { username, password } = this.props;
    api.getServiceStatuses(username, password)
      .then(services => this.setState({ services: services }));
  }

  componentWillUnmount() {
    clearInterval(this.state.intervalId);
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
        <LogoutButton unsetCredentials={this.props.unsetCredentials} />
      </div>
    )
  }

}
