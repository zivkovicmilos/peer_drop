import {
  Avatar,
  Backdrop,
  Box,
  Button,
  Fade,
  IconButton,
  List,
  Modal,
  TextField,
  Typography
} from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';
import CloseRoundedIcon from '@material-ui/icons/CloseRounded';
import { FC, useContext, useEffect, useState } from 'react';
import SessionContext from '../../../context/SessionContext';
import IdentitiesService from '../../../services/identities/identitiesService';
import { IIdentityResponse } from '../../../services/identities/identitiesService.types';
import LoadingIndicator from '../../atoms/LoadingIndicator/LoadingIndicator';
import useSnackbar from '../Snackbar/useSnackbar.hook';
import { ISwapIdentitiesProps } from './swapIdentities.types';

const SwapIdentities: FC<ISwapIdentitiesProps> = (props) => {
  const { modalOpen, setModalOpen } = props;

  const { setUserIdentity } = useContext(SessionContext);

  const [loading, setLoading] = useState<boolean>(true);
  const [allIdentities, setAllidentities] = useState<
    { list: IIdentityResponse[] } | undefined
  >();

  const { openSnackbar } = useSnackbar();

  useEffect(() => {
    if (modalOpen) {
      const fetchAllIdentities = async () => {
        return await IdentitiesService.getIdentities({ page: -1, limit: -1 });
      };

      fetchAllIdentities()
        .then((response) => {
          setAllidentities({ list: response.data });

          setFilteredList({ list: response.data });
        })
        .catch((err) => {
          openSnackbar('Unable to fetch identity list', 'error');
        })
        .finally(() => {
          setLoading(false);
        });
    }
  }, [modalOpen]);

  const [filteredList, setFilteredList] = useState<
    { list: IIdentityResponse[] } | undefined
  >({ list: [] });

  const manualFilter = (fn: any, array: any) => {
    const f = [];

    for (let i = 0; i < array.length; i++) {
      if (fn(array[i])) {
        f.push(array[i]);
      }
    }
    return f;
  };

  const handleSearchInputChange = (event: any) => {
    const value = event.target.value;

    let tempFilteredContacts = { list: allIdentities.list };
    // Filter if there is a value in search input
    if (value) {
      tempFilteredContacts.list = manualFilter(
        (identity: IIdentityResponse) =>
          identity.name.toLowerCase().includes(value.toLowerCase()) ||
          identity.publicKeyID.toLowerCase().includes(value.toLowerCase()),
        allIdentities.list
      );
    }

    setFilteredList(tempFilteredContacts);
  };

  const classes = useStyles();

  const handleIdentitySelect = (identity: IIdentityResponse) => {
    setUserIdentity(identity);
    setModalOpen(false);
  };

  return (
    <Modal
      className={classes.modal}
      open={modalOpen}
      onClose={() => setModalOpen(false)}
      closeAfterTransition
      BackdropComponent={Backdrop}
      BackdropProps={{
        timeout: 500
      }}
    >
      <Fade in={modalOpen}>
        <div className={classes.modalWrapper}>
          <Box
            display={'flex'}
            alignItems={'center'}
            justifyContent={'space-between'}
          >
            <Typography className={classes.modalTitle}>
              Select new identity
            </Typography>
            <IconButton
              classes={{
                root: 'iconButtonRoot'
              }}
              onClick={() => setModalOpen(false)}
            >
              <CloseRoundedIcon
                style={{
                  width: '20px',
                  height: 'auto'
                }}
              />
            </IconButton>
          </Box>
          <Box display={'flex'} flexDirection={'column'} mt={5}>
            <Box mb={4}>
              <TextField
                style={{
                  width: '100%'
                }}
                placeholder={'Search'}
                onChange={handleSearchInputChange}
              />
            </Box>
            <Box
              display={'flex'}
              flexDirection={'column'}
              maxHeight={'350px'}
              overflow={'auto'}
            >
              <List>
                {filteredList &&
                  filteredList.list.map((identity) => {
                    return (
                      <Button
                        classes={{
                          root: classes.buttonOverride
                        }}
                        onClick={() => handleIdentitySelect(identity)}
                      >
                        <Box>
                          <Avatar src={identity.picture}>
                            {identity.name.charAt(0)}
                          </Avatar>
                        </Box>
                        <Box ml={2}>
                          <Typography>{`${identity.name} (${identity.publicKeyID})`}</Typography>
                        </Box>
                      </Button>
                    );
                  })}
                {loading && <LoadingIndicator size={54} />}
                {!loading && filteredList && filteredList.list.length < 1 && (
                  <Box
                    display={'flex'}
                    width={'100%'}
                    justifyContent={'center'}
                    mt={1}
                  >
                    <Typography>No identities found</Typography>
                  </Box>
                )}
              </List>
            </Box>
          </Box>
        </div>
      </Fade>
    </Modal>
  );
};

const useStyles = makeStyles((theme) => {
  return {
    modal: {
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'center'
    },
    modalTitle: {
      fontWeight: 600,
      fontSize: '1.5rem'
    },
    modalWrapper: {
      backgroundColor: theme.palette.background.paper,
      boxShadow: theme.palette.boxShadows.darker,
      padding: '20px 30px',
      // TODO add responsive features
      width: '400px',
      height: '530px',
      border: 'none',
      outline: 'none',
      borderRadius: '15px'
    },
    modalTextMain: {},
    modalSubtext: {
      fontWeight: 500
    },
    buttonOverride: {
      justifyContent: 'left'
    }
  };
});

export default SwapIdentities;
