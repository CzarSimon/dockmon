import React from 'react';
import { Alert } from 'antd';
import './errorMessage.css';

const ErrorMessage = props => (
  (props.text)
    ? <Alert className='error-message' message={props.text} type="error" showIcon closable />
    : <div />);

export default ErrorMessage;
