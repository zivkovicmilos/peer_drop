import { Box, Switch, Typography } from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';
import { FC, useEffect, useState } from 'react';
import theme from '../../../theme/theme';
import { ISwitchBoxProps } from './switchBox.types';

const SwitchBox: FC<ISwitchBoxProps> = (props) => {
  const { type, description, onToggle, toggled, disabled = false } = props;

  const classes = useStyles();

  const [switchToggled, setSwitchToggled] = useState<boolean>(toggled);

  const handleToggle = () => {
    setSwitchToggled(!switchToggled);
  };

  useEffect(() => {
    onToggle(switchToggled);
  }, [switchToggled]);

  return (
    <Box className={classes.switchBoxWrapper}>
      <Box
        display={'flex'}
        flexDirection={'column'}
        mr={2}
        justifyContent={'center'}
      >
        <Typography className={classes.switchBoxTitle}>{type}</Typography>
        <Typography className={classes.switchBoxDescription}>
          {description}
        </Typography>
      </Box>
      <Box
        display={'flex'}
        height={'100%'}
        ml={'auto'}
        justifyContent={'center'}
        alignItems={'center'}
      >
        <Switch
          checked={switchToggled}
          onChange={handleToggle}
          name={type}
          color={'primary'}
          disabled={disabled}
        />
      </Box>
    </Box>
  );
};

const useStyles = makeStyles(() => {
  return {
    switchBoxWrapper: {
      borderRadius: '15px',
      backgroundColor: theme.palette.custom.darkGray,
      height: '150px',
      width: '80%',
      display: 'flex',
      padding: '20px 25px'
    },
    switchBoxTitle: {
      fontWeight: 600,
      marginBottom: '0.5rem'
    },
    switchBoxDescription: {
      color: 'rgba(0, 0, 0, 0.5)',
      textAlign: 'left'
    }
  };
});

export default SwitchBox;
