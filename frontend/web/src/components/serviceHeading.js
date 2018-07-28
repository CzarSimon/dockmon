import React from 'react';
import { Icon } from 'antd';

const ServiceHeading = props => (
  <div className='service-heading'>
    <h3>{props.serviceName}</h3>
    <Icon className='icon'
      type={(props.isHealty) ? 'check-circle' : 'close-circle'}
      style={{fontSize: '18px', color: (props.isHealty) ? '#00BD76' : '#EF472F'}} />
  </div>
)

export default ServiceHeading;
