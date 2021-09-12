import { Box, TextField } from '@material-ui/core';
import { DateTimePicker } from '@material-ui/pickers';
import { useFormik } from 'formik';
import moment from 'moment';
import { FC, useContext, useEffect, useState } from 'react';
import NewWorkspaceContext from '../../../context/NewWorkspaceContext';
import { ContactResponse } from '../../../context/newWorkspaceContext.types';
import { nwPermissionsSchema } from '../../../shared/schemas/workspaceSchemas';
import StepButton from '../../atoms/StepButton/StepButton';
import SwitchBox from '../../atoms/SwitchBox/SwitchBox';
import SpecificContacts from '../SpecificContacts/SpecificContacts';
import {
  ENWPermissionsType,
  INewWorkspacePermissionsProps
} from './newWorkspacePermissions.types';

const NewWorkspacePermissions: FC<INewWorkspacePermissionsProps> = () => {
  const {
    handleNext,
    handleBack,
    permissions,
    setPermissions,
    step,
    accessControlType
  } = useContext(NewWorkspaceContext);

  const [errorMessage, setErrorMessage] = useState<string>('');

  const permissionsFormik = useFormik({
    initialValues: {
      autoCloseWorkspace: permissions.autocloseWorkspace.active,
      autoCloseDate: permissions.autocloseWorkspace.date,

      // Peer limit cannot be enforced if the security
      // type is: Specific contacts
      enforcePeerLimit: permissions.enforcePeerLimit.active,
      peerLimit: permissions.enforcePeerLimit.limit,

      hasAdditionalOwners: permissions.additionalOwners.active,
      contactIDs: permissions.additionalOwners.contactIDs
        ? permissions.additionalOwners.contactIDs
        : []
    },
    validationSchema: nwPermissionsSchema,
    onSubmit: (values, { resetForm }) => {
      if (values.hasAdditionalOwners && values.contactIDs.length < 1) {
        setErrorMessage('At least 1 contact is required');

        return;
      }

      setPermissions({
        autocloseWorkspace: {
          active: values.autoCloseWorkspace,
          date: values.autoCloseWorkspace ? values.autoCloseDate : null
        },
        enforcePeerLimit: {
          active: values.enforcePeerLimit,
          limit: values.peerLimit
        },
        additionalOwners: {
          active: values.hasAdditionalOwners,
          contactIDs: values.contactIDs
        }
      });

      handleNext();
    }
  });

  const [autocloseWorkspace, setAutoCloseWorkspace] = useState<boolean>(
    permissionsFormik.values.autoCloseWorkspace
  );

  const [enforcePeerLimit, setEnforcePeerLimit] = useState<boolean>(
    permissionsFormik.values.enforcePeerLimit
  );

  const [additionalOwners, setAdditionalOwners] = useState<boolean>(
    permissionsFormik.values.hasAdditionalOwners
  );

  const workspacePermissions = [
    {
      type: ENWPermissionsType.AUTO_CLOSE_WORKSPACE,
      description: 'Define when the workspace will close automatically',
      onToggle: (state: boolean) => {
        permissionsFormik.values.autoCloseWorkspace = state;
        setAutoCloseWorkspace(state);
      },
      toggled: permissionsFormik.values.autoCloseWorkspace,
      disabled: false
    },
    {
      type: ENWPermissionsType.ENFORCE_PEER_LIMIT,
      description: 'Define how many users can access the workspace',
      onToggle: (state: boolean) => {
        permissionsFormik.values.enforcePeerLimit = state;

        setEnforcePeerLimit(state);
      },
      toggled: permissionsFormik.values.enforcePeerLimit,
      // disabled: accessControlType == ENWAccessControl.SPECIFIC_CONTACTS TODO bring back
      disabled: false
    },
    {
      type: ENWPermissionsType.ADDITIONAL_WORKSPACE_OWNERS,
      description: 'Define who else can control the workspace',
      onToggle: (state: boolean) => {
        permissionsFormik.values.hasAdditionalOwners = state;

        setAdditionalOwners(state);
      },
      toggled: permissionsFormik.values.hasAdditionalOwners,
      disabled: false
    }
  ];

  // Renders a single switch box based on the ordering
  const renderSwitchBox = (index: number) => {
    const switchBox = workspacePermissions[index];

    return (
      <SwitchBox
        key={`switch-${index}`}
        type={switchBox.type}
        description={switchBox.description}
        onToggle={switchBox.onToggle}
        toggled={switchBox.toggled}
        disabled={switchBox.disabled}
      />
    );
  };

  const getDate = () => {
    if (permissionsFormik.values.autoCloseDate != null) {
      return moment(permissionsFormik.values.autoCloseDate).toDate();
    }

    return moment().toDate();
  };

  const [selectedDate, setSelectedDate] = useState<Date>(getDate);

  const [contactIDs, setContactIDs] = useState<{ contacts: ContactResponse[] }>(
    { contacts: permissionsFormik.values.contactIDs }
  );

  useEffect(() => {
    permissionsFormik.values.contactIDs = contactIDs.contacts;
  }, [contactIDs]);

  return (
    <Box display={'flex'} flexDirection={'column'} width={'100%'}>
      <form autoComplete={'off'} onSubmit={permissionsFormik.handleSubmit}>
        {/*Top row*/}
        <Box display={'flex'} width={'100%'}>
          <Box display={'flex'} flexDirection={'column'} width={'50%'}>
            {renderSwitchBox(0)}
            <Box mt={4}>
              <DateTimePicker
                autoOk
                disabled={!autocloseWorkspace}
                disablePast
                ampm={false}
                value={selectedDate}
                onChange={(date) => setSelectedDate(moment(date).toDate())}
                label={'Close date/time'}
              />
            </Box>
          </Box>
          <Box display={'flex'} flexDirection={'column'} width={'50%'}>
            {renderSwitchBox(1)}
            <Box minHeight={'80px'} mt={4} width={'50%'}>
              <TextField
                id={'peerLimit'}
                label={'Peer limit'}
                variant={'outlined'}
                value={permissionsFormik.values.peerLimit}
                onChange={permissionsFormik.handleChange}
                onBlur={permissionsFormik.handleBlur}
                error={
                  permissionsFormik.touched.peerLimit &&
                  Boolean(permissionsFormik.errors.peerLimit)
                }
                helperText={
                  permissionsFormik.touched.peerLimit &&
                  permissionsFormik.errors.peerLimit
                }
                disabled={!enforcePeerLimit}
              />
            </Box>
          </Box>
        </Box>
        {/*Bottom row*/}
        <Box display={'flex'} width={'100%'} mt={4}>
          <Box display={'flex'} width={'50%'}>
            <SpecificContacts
              contactIDs={contactIDs}
              setContactIDs={setContactIDs}
              errorMessage={errorMessage}
              listTitle={'Additional owners'}
              disabled={!additionalOwners}
            />
          </Box>
          <Box display={'flex'} width={'50%'}>
            {renderSwitchBox(2)}
          </Box>
        </Box>
        <Box display={'flex'} alignItems={'center'} mt={4}>
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

export default NewWorkspacePermissions;
