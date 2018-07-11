import React, { Component } from 'react';

export default class ServiceStatus extends Component {

  render() {
    const { serviceName, isHealty } = this.props;
    console.log(this.props);
    return (
      <div>
        <p>{serviceName}</p>
        {(isHealty) ? "Healty" : "Unhealthy"}
      </div>
    )
  }

}
