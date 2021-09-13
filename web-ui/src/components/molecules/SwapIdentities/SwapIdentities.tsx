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
import { IUserIdentity } from '../../../context/sessionContext.types';
import LoadingIndicator from '../../atoms/LoadingIndicator/LoadingIndicator';
import useSnackbar from '../Snackbar/useSnackbar.hook';
import { ISwapIdentitiesProps } from './swapIdentities.types';

const SwapIdentities: FC<ISwapIdentitiesProps> = (props) => {
  const { modalOpen, setModalOpen } = props;

  const { setUserIdentity } = useContext(SessionContext);

  const [loading, setLoading] = useState<boolean>(true);
  const [allIdentities, setAllidentities] = useState<
    { list: IUserIdentity[] } | undefined
  >();

  const { openSnackbar } = useSnackbar();

  useEffect(() => {
    if (modalOpen) {
      const fetchAllContacts = async () => {
        return [
          {
            picture: '123',
            name: 'Milos Zivkovic',
            keyID: '123'
          },
          {
            picture: '1234',
            name: 'Milos Zivkovic',
            keyID: '1234'
          },
          {
            picture: '12345',
            name: 'Milos Zivkovic',
            keyID: '12345'
          },
          {
            picture: '123456',
            name: 'Milos Zivkovic',
            keyID: '123456'
          },
          {
            picture: '1234567',
            name: 'Milos Zivkovic',
            keyID: '123456'
          },
          {
            picture: '12345678',
            name: 'Milos Zivkovic',
            keyID: '123456'
          },
          {
            picture: '12345689',
            name: 'Milos Zivkovic',
            keyID: '123456'
          },
          {
            picture: '1234568910',
            name: 'Milos Zivkovic',
            keyID: '123456'
          },
          {
            picture: '1234568911',
            name: 'Milos Zivkovic',
            keyID: '123456'
          }
        ];
      };

      fetchAllContacts()
        .then((response) => {
          setAllidentities({ list: response });

          setFilteredList({ list: response });
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
    { list: IUserIdentity[] } | undefined
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
        (contact: any) =>
          contact.name.toLowerCase().includes(value.toLowerCase()) ||
          contact.keyID.toLowerCase().includes(value.toLowerCase()),
        allIdentities.list
      );
    }

    setFilteredList(tempFilteredContacts);
  };

  const classes = useStyles();

  const handleIdentitySelect = (identity: IUserIdentity) => {
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
                          <Typography>{`${identity.name} (${identity.keyID})`}</Typography>
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
