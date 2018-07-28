import React, { Component } from 'react';
import { Badge } from 'antd';
import { date } from '../utils';

export default class ServiceInfo extends Component {

  render() {
    const {
      shouldRestart, restarts, lastRestarted,
      lastHealthFailure, lastHealthSuccess, createdAt
    } = this.props;
    return (
      <div className='service-info'>
        <p>Restarts: {(restarts) ? <Badge count={restarts} /> : restarts}</p>
        <p>Last Restarted: {date.fromNow(lastRestarted)}</p>
        <p>Restarts on failure: {(shouldRestart) ? 'Yes' : 'No'}</p>
        <br />
        <p>Last healthy: {date.fromNow(lastHealthSuccess)}</p>
        <p>Last unhealthy: {date.fromNow(lastHealthFailure)}</p>
        <br />
        <p>Tracking started: {date.toDateString(createdAt)}</p>
      </div>
    )
  }

}
