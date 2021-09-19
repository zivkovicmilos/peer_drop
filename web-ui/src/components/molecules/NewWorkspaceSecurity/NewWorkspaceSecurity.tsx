import {
  Box,
  FormControlLabel,
  Radio,
  RadioGroup,
  Tooltip
} from '@material-ui/core';
import InfoRoundedIcon from '@material-ui/icons/InfoRounded';
import { useFormik } from 'formik';
import { FC, useContext, useEffect, useState } from 'react';
import NewWorkspaceContext from '../../../context/NewWorkspaceContext';
import {
  ENWAccessControl,
  INWAccessControlContacts,
  INWAccessControlPassword
} from '../../../context/newWorkspaceContext.types';
import { IContactResponse } from '../../../services/contacts/contactsService.types';
import { nwSecurityPasswordSchema } from '../../../shared/schemas/workspaceSchemas';
import theme from '../../../theme/theme';
import FormTitle from '../../atoms/FormTitle/FormTitle';
import StepButton from '../../atoms/StepButton/StepButton';
import PasswordStrength from '../PasswordStrength/PasswordStrength';
import SpecificContacts from '../SpecificContacts/SpecificContacts';
import { INewWorkspaceSecurityProps } from './newWorkspaceSecurity.types';

const NewWorkspaceSecurity: FC<INewWorkspaceSecurityProps> = () => {
  const {
    handleNext,
    handleBack,
    accessControl,
    accessControlType,
    setAccessControlType,
    setAccessControl,
    step
  } = useContext(NewWorkspaceContext);

  const [errorMessage, setErrorMessage] = useState<string>('');

  const securityFormik = useFormik({
    initialValues: {
      accessControlType: accessControlType,
      password:
        accessControlType == ENWAccessControl.PASSWORD
          ? (accessControl as INWAccessControlPassword).password
          : '',
      passwordConfirm: ENWAccessControl.PASSWORD
        ? (accessControl as INWAccessControlPassword).password
        : '',
      contactIDs:
        accessControlType == ENWAccessControl.SPECIFIC_CONTACTS
          ? (accessControl as INWAccessControlContacts).contacts
          : []
    },
    validationSchema: nwSecurityPasswordSchema,
    onSubmit: (values, { resetForm }) => {
      setErrorMessage('');

      if (accessControlType == ENWAccessControl.PASSWORD) {
        setAccessControl({ password: values.password });
      } else {
        if (values.contactIDs.length < 1) {
          setErrorMessage('At least 1 contact is required');

          return;
        }

        setAccessControl({ contacts: values.contactIDs });
      }

      handleNext();
    }
  });

  const securityTypes = [
    ENWAccessControl.SPECIFIC_CONTACTS,
    ENWAccessControl.PASSWORD
  ];

  const handleSecurityTypeSelect = (type: ENWAccessControl) => {
    setAccessControlType(type);
    securityFormik.values.accessControlType = type;
  };

  const [contactIDs, setContactIDs] = useState<{
    contacts: IContactResponse[];
  }>({ contacts: securityFormik.values.contactIDs });

  useEffect(() => {
    securityFormik.values.contactIDs = contactIDs.contacts;
  }, [contactIDs]);

  return (
    <Box display={'flex'} flexDirection={'column'} width={'50%'}>
      <form autoComplete={'off'} onSubmit={securityFormik.handleSubmit}>
        <Box display={'flex'} flexDirection={'column'} maxWidth={'100%'} mb={4}>
          <Box mb={2} display={'flex'} alignItems={'center'}>
            <FormTitle title={'Access control'} />
            <Box ml={0.5}>
              <Tooltip
                title="Level of access to the workspace"
                arrow
                placement={'top'}
              >
                <InfoRoundedIcon
                  style={{
                    width: '15px',
                    height: 'auto',
                    color: theme.palette.custom.transparentBlack
                  }}
                />
              </Tooltip>
            </Box>
          </Box>
          <RadioGroup name={'accessControlType'} value={accessControlType} row>
            {securityTypes.map((securityType) => {
              return (
                <FormControlLabel
                  value={securityType}
                  style={{
                    marginLeft: 0
                  }}
                  control={<Radio color={'primary'} />}
                  label={securityType}
                  onChange={() => {
                    handleSecurityTypeSelect(securityType);
                  }}
                />
              );
            })}
          </RadioGroup>
        </Box>
        {accessControlType == ENWAccessControl.PASSWORD && (
          <Box mb={8}>
            <Box mb={2}>
              <FormTitle title={'Set password'} />
            </Box>
            <PasswordStrength formik={securityFormik} />
          </Box>
        )}

        {accessControlType == ENWAccessControl.SPECIFIC_CONTACTS && (
          <Box mb={4}>
            <SpecificContacts
              contactIDs={contactIDs}
              setContactIDs={setContactIDs}
              errorMessage={errorMessage}
            />
          </Box>
        )}
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

export default NewWorkspaceSecurity;
