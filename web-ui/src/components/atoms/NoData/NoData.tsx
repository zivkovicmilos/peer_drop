import { Box, Typography } from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';
import { FC } from 'react';
import { ReactComponent as Clipboard } from '../../../shared/assets/icons/clipboard.svg';
import { INoDataProps } from './noData.types';

const NoData: FC<INoDataProps> = (props) => {
  const { text } = props;

  const classes = useStyles();

  return (
    <Box className={classes.noDataWrapper}>
      <Clipboard className={classes.clipboard} />
      <Box mt={3}>
        <Typography>{text}</Typography>
      </Box>
    </Box>
  );
};

const useStyles = makeStyles(() => {
  return {
    noDataWrapper: {
      display: 'flex',
      marginLeft: 'auto',
      marginRight: 'auto',
      flexDirection: 'column',
      justifyContent: 'center',
      alignItems: 'center',
      opacity: 0.5
    },
    clipboard: {
      width: '160px',
      height: 'auto'
    }
  };
});

export default NoData;
