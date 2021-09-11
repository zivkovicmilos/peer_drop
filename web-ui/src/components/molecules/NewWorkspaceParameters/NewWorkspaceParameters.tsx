import {
  Box,
  FormControlLabel,
  Radio,
  RadioGroup,
  TextField,
  Typography
} from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';
import MoveToInboxRoundedIcon from '@material-ui/icons/MoveToInboxRounded';
import SendRoundedIcon from '@material-ui/icons/SendRounded';
import SyncAltRoundedIcon from '@material-ui/icons/SyncAltRounded';
import { useFormik } from 'formik';
import { FC, useContext } from 'react';
import NewWorkspaceContext from '../../../context/NewWorkspaceContext';
import { ENewWorkspaceType } from '../../../context/newWorkspaceContext.types';
import { nwParametersSchema } from '../../../shared/schemas/workspaceSchemas';
import theme from '../../../theme/theme';
import FormTitle from '../../atoms/FormTitle/FormTitle';
import StepButton from '../../atoms/StepButton/StepButton';
import { INewWorkspaceParametersProps } from './newWorkspace.types';

const NewWorkspaceParameters: FC<INewWorkspaceParametersProps> = (props) => {
  const {
    workspaceName,
    workspaceType,
    handleNext,
    handleBack,
    setWorkspaceType,
    setWorkspaceName,
    step
  } = useContext(NewWorkspaceContext);

  const parametersFormik = useFormik({
    initialValues: {
      workspaceName,
      workspaceType
    },
    validationSchema: nwParametersSchema,
    onSubmit: (values, { resetForm }) => {
      setWorkspaceName(values.workspaceName);

      handleNext();
    }
  });

  const workspaceTypes = [
    {
      type: ENewWorkspaceType.SEND_ONLY,
      description:
        'The workspace owner can distribute files, and workspace participants can only receive them.'
    },
    {
      type: ENewWorkspaceType.RECEIVE_ONLY,
      description:
        'The workspace owner can only receive files, and workspace participants can only distribute them.'
    },
    {
      type: ENewWorkspaceType.SEND_RECEIVE,
      description: 'All workspace participants can send and receive files. '
    }
  ];

  const handleTypeSelect = (type: ENewWorkspaceType) => {
    setWorkspaceType(type);
  };

  const classes = useStyles();

  const renderTypeIcon = (type: ENewWorkspaceType) => {
    switch (type) {
      case ENewWorkspaceType.SEND_ONLY:
        return <SendRoundedIcon />;
      case ENewWorkspaceType.RECEIVE_ONLY:
        return <MoveToInboxRoundedIcon />;
      case ENewWorkspaceType.SEND_RECEIVE:
        return <SyncAltRoundedIcon />;
    }
  };

  return (
    <Box display={'flex'} flexDirection={'column'} width={'50%'}>
      <form autoComplete={'off'} onSubmit={parametersFormik.handleSubmit}>
        <Box display={'flex'} flexDirection={'column'} maxWidth={'50%'} mb={2}>
          <Box mb={2}>
            <FormTitle title={'Basic info'} />
          </Box>
          <Box minHeight={'80px'}>
            <TextField
              id={'workspaceName'}
              label={'Name'}
              variant={'outlined'}
              value={parametersFormik.values.workspaceName}
              onChange={parametersFormik.handleChange}
              onBlur={parametersFormik.handleBlur}
              error={
                parametersFormik.touched.workspaceName &&
                Boolean(parametersFormik.errors.workspaceName)
              }
              helperText={
                parametersFormik.touched.workspaceName &&
                parametersFormik.errors.workspaceName
              }
            />
          </Box>
        </Box>
        <Box display={'flex'} flexDirection={'column'} mb={2}>
          <Box mb={2}>
            <FormTitle title={'Workspace type'} />
          </Box>
          <RadioGroup name={'workspaceType'} value={workspaceType}>
            {workspaceTypes.map((typeWrapper, index) => {
              return (
                <Box display={'flex'} flexDirection={'column'}>
                  <FormControlLabel
                    value={typeWrapper.type}
                    style={{
                      marginLeft: 0
                    }}
                    control={<Radio color={'primary'} />}
                    label={
                      <Box display={'flex'} alignItems={'center'}>
                        <Typography className={classes.workspaceType}>
                          {typeWrapper.type}
                        </Typography>
                        <Box ml={1} display={'flex'} alignItems={'center'}>
                          {renderTypeIcon(typeWrapper.type)}
                        </Box>
                      </Box>
                    }
                    onChange={() => {
                      handleTypeSelect(typeWrapper.type);
                    }}
                  />
                  <Box ml={'30px'} mb={2}>
                    <Typography className={classes.workspaceTypeDescription}>
                      {typeWrapper.description}
                    </Typography>
                  </Box>
                </Box>
              );
            })}
          </RadioGroup>
        </Box>
        <Box display={'flex'} alignItems={'center'}>
          <Box mr={2}>
            <StepButton
              text={'Back'}
              variant={'outlined'}
              disabled={step < 1}
              shouldSubmit={false}
              onClick={() => {
                handleBack();
              }}
            />
          </Box>
          <StepButton text={'Next'} variant={'contained'} />
        </Box>
      </form>
    </Box>
  );
};

const useStyles = makeStyles(() => {
  return {
    workspaceType: {
      fontWeight: 600
    },
    workspaceTypeDescription: {
      color: theme.palette.custom.transparentBlack
    }
  };
});

export default NewWorkspaceParameters;
