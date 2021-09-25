import {
  Avatar,
  Box,
  Button,
  ClickAwayListener,
  Grow,
  InputBase,
  Paper,
  Popper,
  Typography
} from '@material-ui/core';
import { makeStyles, Theme } from '@material-ui/core/styles';
import SearchRoundedIcon from '@material-ui/icons/SearchRounded';
import React, { FC, useContext, useRef, useState } from 'react';
import { useHistory } from 'react-router-dom';
import SessionContext from '../../../context/SessionContext';
import { IContactResponse } from '../../../services/contacts/contactsService.types';
import { IIdentityResponse } from '../../../services/identities/identitiesService.types';
import SearchService from '../../../services/search/searchService';
import { IWorkspaceWrapper } from '../../../services/workspaces/workspacesService.types';
import CommonUtils from '../../../shared/utils/CommonUtils';
import theme from '../../../theme/theme';
import FormTitle from '../../atoms/FormTitle/FormTitle';
import LoadingIndicator from '../../atoms/LoadingIndicator/LoadingIndicator';
import useSnackbar from '../Snackbar/useSnackbar.hook';
import { ISearchBarProps } from './searchbar.types';

const Searchbar: FC<ISearchBarProps> = () => {
  const classes = useStyles();
  const { searchContext } = useContext(SessionContext);

  const anchorRef = useRef<HTMLDivElement>(null);

  const [open, setOpen] = useState<boolean>(false);
  const id = open ? 'search-popover' : undefined;

  const handleClose = (event: any) => {
    if (
      anchorRef.current &&
      anchorRef.current.contains(event.target as HTMLElement)
    ) {
      return;
    }

    setOpen(false);
  };

  const history = useHistory();

  const handleWorkspaceClick = (mnemonic: string) => {
    setOpen(false);

    history.push(`/workspaces/view/${CommonUtils.formatMnemonic(mnemonic)}`);
  };

  const renderWorkspacesSection = (workspaces: IWorkspaceWrapper[]) => {
    if (workspaces && workspaces.length > 0) {
      return (
        <Box display={'flex'} flexDirection={'column'}>
          <Box>
            <FormTitle title={'Workspaces'} />
          </Box>
          <Box display={'flex'} flexDirection={'column'} width={'100%'}>
            {workspaces.map((result) => {
              return (
                <Button
                  key={result.workspaceMnemonic}
                  classes={{
                    root: classes.buttonOverride
                  }}
                  onClick={() => {
                    handleWorkspaceClick(result.workspaceName);
                  }}
                >
                  <Box
                    display={'flex'}
                    padding={'10px 15px'}
                    alignItems={'center'}
                  >
                    <Box>{result.workspaceName}</Box>
                  </Box>
                </Button>
              );
            })}
          </Box>
        </Box>
      );
    }
  };

  const handleContactClick = (id: string) => {
    setOpen(false);

    history.push(`/contacts/${id}/edit`);
  };

  const renderContactsSection = (contacts: IContactResponse[]) => {
    if (contacts && contacts.length > 0) {
      return (
        <Box display={'flex'} flexDirection={'column'}>
          <Box>
            <FormTitle title={'Contacts'} />
          </Box>
          <Box display={'flex'} flexDirection={'column'} width={'100%'}>
            {contacts.map((result) => {
              return (
                <Button
                  key={result.id}
                  classes={{
                    root: classes.buttonOverride
                  }}
                  onClick={() => {
                    handleContactClick(result.id);
                  }}
                >
                  <Box
                    display={'flex'}
                    padding={'10px 15px'}
                    alignItems={'center'}
                  >
                    <Box>{result.name}</Box>
                  </Box>
                </Button>
              );
            })}
          </Box>
        </Box>
      );
    }
  };

  const handleIdentityNav = (id: string) => {
    setOpen(false);

    history.push(`/identities/${id}/edit`);
  };

  const renderIdentitiesSection = (identities: IIdentityResponse[]) => {
    if (identities && identities.length > 0) {
      return (
        <Box display={'flex'} flexDirection={'column'}>
          <Box>
            <FormTitle title={'Identities'} />
          </Box>
          <Box display={'flex'} flexDirection={'column'} width={'100%'}>
            {identities.map((result) => {
              return (
                <Button
                  key={result.id}
                  classes={{
                    root: classes.buttonOverride
                  }}
                  onClick={() => {
                    handleIdentityNav(result.id);
                  }}
                >
                  <Box
                    display={'flex'}
                    padding={'10px 15px'}
                    alignItems={'center'}
                  >
                    <Avatar src={result.picture} />
                    <Box ml={2}>{result.name}</Box>
                  </Box>
                </Button>
              );
            })}
          </Box>
        </Box>
      );
    }
  };

  const handleToggle = () => {
    setOpen(!open);
  };

  const [searchItems, setSearchItems] = useState<any>(null);

  const fetchSearchResults = async (input: string) => {
    return await SearchService.getSearchResults(input);
  };

  const [searchValue, setSearchValue] = useState<string>('');

  const { openSnackbar } = useSnackbar();

  const handleSearch = () => {
    if (searchValue) {
      setOpen(true);

      fetchSearchResults(searchValue)
        .then((response) => {
          let items = [];
          let identities = renderIdentitiesSection(response.identities);
          if (identities) {
            items.push(identities);
          }

          let workspaces = renderWorkspacesSection(response.workspaces);
          if (workspaces) {
            items.push(workspaces);
          }

          let contacts = renderContactsSection(response.contacts);
          if (contacts) {
            items.push(contacts);
          }

          setSearchItems(items);
        })
        .catch((err) => {
          openSnackbar('Unable to perform search', 'error');
        });
    } else {
      setOpen(false);
    }
  };

  const renderSearchItems = () => {
    if (searchItems && searchItems.length > 0) {
      return searchItems;
    } else if (searchItems && searchItems.length == 0) {
      return (
        <Typography>
          <i>No search results</i>
        </Typography>
      );
    } else {
      return <LoadingIndicator size={20} />;
    }
  };

  return (
    <Paper component={'form'} className={classes.searchBar} ref={anchorRef}>
      <SearchRoundedIcon
        style={{
          fill: theme.palette.custom.transparentBlack
        }}
      />

      <InputBase
        className={classes.input}
        placeholder={'Search'}
        onChange={(event) => {
          setSearchValue(event.target.value);
        }}
        onKeyUp={() => {
          handleSearch();
        }}
      />

      <Popper
        open={open}
        anchorEl={anchorRef.current}
        role={undefined}
        transition
        placement={'bottom-start'}
      >
        {({ TransitionProps, placement }) => (
          <Grow {...TransitionProps}>
            <ClickAwayListener onClickAway={handleClose}>
              <Box
                width={'400px'}
                style={{
                  backgroundColor: theme.palette.custom.darkGray,
                  padding: '10px 15px',
                  borderRadius: '15px'
                }}
              >
                {renderSearchItems()}
              </Box>
            </ClickAwayListener>
          </Grow>
        )}
      </Popper>
    </Paper>
  );
};

const useStyles = makeStyles((theme: Theme) => {
  return {
    searchBar: {
      display: 'flex',
      minWidth: '75%', // TODO change this to be more dynamic
      alignItems: 'center',
      borderRadius: '10px',
      padding: '10px 15px',
      backgroundColor: theme.palette.custom.darkGray,
      boxShadow: 'none',
      outline: 'none',
      '&:focus-within': {
        boxShadow: theme.palette.boxShadows.darker
      }
    },
    input: {
      marginLeft: theme.spacing(1),
      flex: 1
    },
    buttonOverride: {
      justifyContent: 'left'
    }
  };
});

export default Searchbar;
