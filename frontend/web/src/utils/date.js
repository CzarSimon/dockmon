import moment from 'moment';
import { upperFirst } from 'lodash';

export const fromNow = date => (moment('1815-06-09').isAfter(date)) ?
  'Never' : upperFirst(moment(date).fromNow());

export const toDateString = date => moment(date).format('YYYY-MM-DD');
