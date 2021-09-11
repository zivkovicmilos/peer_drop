import { Box, TextField, Typography } from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';
import { passwordStrength } from 'check-password-strength';
import { FC, useEffect, useState } from 'react';
import { ReactComponent as PasswordBar } from '../../../shared/assets/icons/password_bar.svg';
import theme from '../../../theme/theme';
import {
  EPasswordStrength,
  IPasswordStrengthProps
} from './passwordStrength.types';

const PasswordStrength: FC<IPasswordStrengthProps> = (props) => {
  const { formik } = props;

  const [psLevel, setPSLevel] = useState<EPasswordStrength>(
    EPasswordStrength.NULL
  );

  useEffect(() => {
    handlePasswordChange();
  }, [formik.values.password]);

  const handlePasswordChange = () => {
    if (formik.values.password == '' || formik.values.password == undefined) {
      setPSLevel(EPasswordStrength.NULL);
      return;
    }

    let ps = passwordStrength(formik.values.password);

    switch (ps.id) {
      case 0:
        setPSLevel(EPasswordStrength.TOO_WEAK);
        return;
      case 1:
        setPSLevel(EPasswordStrength.WEAK);
        return;
      case 2:
        setPSLevel(EPasswordStrength.MEDIUM);
        return;
      case 3:
        setPSLevel(EPasswordStrength.STRONG);
        return;
    }
  };

  const renderGrayBar = (index: number) => {
    return (
      <PasswordBar
        key={`bar-${index}`}
        style={{
          height: '6px',
          width: 'auto',
          marginLeft: '8px',
          fill: 'rgba(0, 0, 0, 0.2)',
          transition: 'fill .24s ease-in-out'
        }}
      />
    );
  };

  const renderGreenBar = (index: number) => {
    return (
      <PasswordBar
        key={`bar-${index}`}
        style={{
          height: '6px',
          width: 'auto',
          marginLeft: '8px',
          fill: '#B5E48C',
          transition: 'fill .24s ease-in-out'
        }}
      />
    );
  };

  const renderPasswordBars = () => {
    let passwordBars = [];

    for (let i = 1; i < 5; i++) {
      if (psLevel == EPasswordStrength.NULL) {
        passwordBars.push(renderGrayBar(i));
        continue;
      }

      if (psLevel.valueOf() >= i) {
        passwordBars.push(renderGreenBar(i));
      } else {
        passwordBars.push(renderGrayBar(i));
      }
    }

    return passwordBars;
  };

  const classes = useStyles();

  return (
    <Box display={'flex'} flexDirection={'column'} width={'100%'}>
      <Box display={'flex'} width={'100%'} mb={2}>
        <Box minHeight={'80px'}>
          <TextField
            id={'password'}
            type={'password'}
            label={'Password'}
            variant={'outlined'}
            value={formik.values.password}
            onChange={formik.handleChange}
            onBlur={formik.handleBlur}
            error={formik.touched.password && Boolean(formik.errors.password)}
            helperText={formik.touched.password && formik.errors.password}
          />
        </Box>
        <Box minHeight={'80px'} ml={4}>
          <TextField
            id={'passwordConfirm'}
            type={'password'}
            label={'Confirm password'}
            variant={'outlined'}
            value={formik.values.passwordConfirm}
            onChange={formik.handleChange}
            onBlur={formik.handleBlur}
            error={
              formik.touched.passwordConfirm &&
              Boolean(formik.errors.passwordConfirm)
            }
            helperText={
              formik.touched.passwordConfirm && formik.errors.passwordConfirm
            }
          />
        </Box>
      </Box>
      <Box display={'flex'} justifyContent={'center'} flexDirection={'column'}>
        <Box display={'flex'} ml={'-8px'} mb={1}>
          {renderPasswordBars()}
        </Box>
        <Typography className={classes.passwordStrengthText}>
          Password strength
        </Typography>
      </Box>
    </Box>
  );
};

const useStyles = makeStyles(() => {
  return {
    passwordStrengthText: {
      fontSize: theme.typography.pxToRem(12),
      fontWeight: 400
    }
  };
});

export default PasswordStrength;
