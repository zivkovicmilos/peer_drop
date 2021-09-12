import { Box, Typography } from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';
import { FC, Fragment } from 'react';
import {
  ENWAccessControl,
  INWAccessControlContacts
} from '../../../context/newWorkspaceContext.types';
import theme from '../../../theme/theme';
import { INewWorkspaceInfoProps } from './newWorkspaceInfo.types';

const NewWorkspaceInfo: FC<INewWorkspaceInfoProps> = (props) => {
  const {
    workspaceName,
    workspaceType,
    accessControl,
    permissions,
    accessControlType
  } = props;

  const classes = useStyles();

  const renderSecurityInformation = () => {
    if (accessControlType == ENWAccessControl.SPECIFIC_CONTACTS) {
      return `${accessControlType} (${
        (accessControl as INWAccessControlContacts).contacts.length
      })`;
    } else {
      return accessControlType;
    }
  };

  const renderAdditionalSettings = () => {
    let additionalSettings = [];

    if (permissions.additionalOwners.active) {
      additionalSettings.push(
        <Typography key={'additional-owners'} className={classes.reviewItem}>
          {`Additional workspace owners - ${permissions.additionalOwners.contactIDs.length}`}
        </Typography>
      );
    }

    if (permissions.autocloseWorkspace.active) {
      additionalSettings.push(
        <Typography key={'autoclose'} className={classes.reviewItem}>
          {`Auto-close workspace - ${permissions.autocloseWorkspace.date}`}
        </Typography>
      );
    }

    if (permissions.enforcePeerLimit.active) {
      additionalSettings.push(
        <Typography key={'peer-limit'} className={classes.reviewItem}>
          {`Enforce peer limit - ${permissions.enforcePeerLimit.limit}`}
        </Typography>
      );
    }

    if (additionalSettings.length < 1) {
      additionalSettings.push(
        <Typography key={'no-settings'} className={classes.reviewItem}>
          No additional settings
        </Typography>
      );
    }

    return additionalSettings;
  };

  return (
    <Fragment>
      <Box display={'flex'} width={'100%'} justifyContent={'space-between'}>
        <Box display={'flex'} flexDirection={'column'}>
          <Typography className={classes.reviewTitle}>Name</Typography>
          <Typography className={classes.reviewItem}>
            {workspaceName}
          </Typography>
        </Box>
        <Box display={'flex'} flexDirection={'column'}>
          <Typography className={classes.reviewTitle}>Type</Typography>
          <Typography className={classes.reviewItem}>
            {workspaceType}
          </Typography>
        </Box>
        <Box display={'flex'} flexDirection={'column'}>
          <Typography className={classes.reviewTitle}>Security</Typography>
          <Typography className={classes.reviewItem}>
            {renderSecurityInformation()}
          </Typography>
        </Box>
      </Box>
      <Box display={'flex'} flexDirection={'column'} mt={6} mb={8}>
        <Typography className={classes.reviewTitle}>
          Additional settings
        </Typography>
        {renderAdditionalSettings()}
      </Box>
    </Fragment>
  );
};

const useStyles = makeStyles(() => {
  return {
    reviewTitle: {
      fontWeight: 600,
      fontSize: theme.typography.pxToRem(16)
    },
    reviewItem: {
      fontWeight: 600,
      color: '#9C9C9C',
      fontSize: theme.typography.pxToRem(16)
    }
  };
});

export default NewWorkspaceInfo;
