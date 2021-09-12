import {
  Button,
  ButtonGroup,
  ClickAwayListener,
  Grow,
  MenuItem,
  MenuList,
  Paper,
  Popper
} from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';
import KeyboardArrowDownRoundedIcon from '@material-ui/icons/KeyboardArrowDownRounded';
import KeyboardArrowUpRoundedIcon from '@material-ui/icons/KeyboardArrowUpRounded';
import clsx from 'clsx';
import React, { FC, Fragment, useRef, useState } from 'react';
import { useHistory } from 'react-router-dom';
import Link from '../../atoms/Link/Link';
import { IMenuActionButtonProps } from './menuActionButton.types';

const MenuActionButton: FC<IMenuActionButtonProps> = () => {
  const [open, setOpen] = useState(false);
  const anchorRef = useRef<HTMLDivElement>(null);

  const classes = useStyles();
  const history = useHistory();

  const handleToggle = () => {
    setOpen((prevOpen) => !prevOpen);
  };

  const handleClose = (event: any) => {
    if (
      anchorRef.current &&
      anchorRef.current.contains(event.target as HTMLElement)
    ) {
      return;
    }

    setOpen(false);
  };

  return (
    <Fragment>
      <ButtonGroup
        variant="contained"
        ref={anchorRef}
        style={{
          border: 'none',
          outline: 'none',
          boxShadow: 'none',
          borderRadius: '15px'
        }}
      >
        <Button
          onClick={() => {
            history.push('/workspaces/new');
          }}
          className={clsx(classes.mainButton, classes.rounded)}
          color={'primary'}
        >
          New Workspace
        </Button>
        <Button
          color={'primary'}
          size={'small'}
          onClick={handleToggle}
          className={classes.rounded}
        >
          {!open ? (
            <KeyboardArrowDownRoundedIcon />
          ) : (
            <KeyboardArrowUpRoundedIcon />
          )}
        </Button>
      </ButtonGroup>
      <Popper
        open={open}
        anchorEl={anchorRef.current}
        role={undefined}
        transition
        disablePortal
      >
        {({ TransitionProps, placement }) => (
          <Grow
            {...TransitionProps}
            style={{
              transformOrigin:
                placement === 'bottom' ? 'center top' : 'center bottom'
            }}
          >
            <Paper>
              <ClickAwayListener onClickAway={handleClose}>
                <Link
                  to={'/workspaces/join'}
                  style={{ textDecoration: 'none' }}
                  onClick={(event) => {
                    handleClose(event);
                  }}
                >
                  <MenuList
                    style={{
                      background: '#303030',
                      borderRadius: '5px',
                      boxShadow: 'none'
                    }}
                  >
                    <MenuItem className={classes.subButton}>
                      Join Workspace
                    </MenuItem>
                  </MenuList>
                </Link>
              </ClickAwayListener>
            </Paper>
          </Grow>
        )}
      </Popper>
    </Fragment>
  );
};

const useStyles = makeStyles(() => {
  return {
    subButton: {
      fontFamily: 'Montserrat',
      fontWeight: 500,
      fontSize: '0.875rem',
      color: 'white'
    },
    rounded: {
      borderRadius: '15px'
    },
    mainButton: {
      padding: '16px 25px',
      fontFamily: 'Montserrat',
      fontWeight: 500,
      borderRadius: '15px',
      border: 'none',
      outline: 'none'
    }
  };
});

export default MenuActionButton;
