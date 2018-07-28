export const LANDSCAPE_ORIENTATION = 'landscape'
export const PORTRAIT_ORIENTATION = 'portrait'

export const getScreenOrientation = () => {
  const { innerHeight, innerWidth } = window;
  return (innerHeight > 1.1 * innerWidth) ? PORTRAIT_ORIENTATION : LANDSCAPE_ORIENTATION;
}

export const getStyle = () => orientationStyles[getScreenOrientation()]

const orientationStyles = {
  [LANDSCAPE_ORIENTATION]: {
    marginLeft: '20%',
    marginRight: '20%',
    width: '60%',
  },
  [PORTRAIT_ORIENTATION]: {
    marginLeft: '5%',
    marginRight: '5%',
    width: '90%',
  }
}
