import { Box } from '@material-ui/core';
import { FC } from 'react';
import { CSSTransition, TransitionGroup } from 'react-transition-group';
import SuggestedFile from '../../atoms/SuggestedFile/SuggestedFile';
import { ISuggestedListProps } from './suggestedList.types';

const SuggestedList: FC<ISuggestedListProps> = (props) => {
  const { files, workspaceMnemonic } = props;

  return (
    <Box ml={-4} display={'flex'} width={'80%'}>
      <TransitionGroup
        style={{
          display: 'flex',
          width: '80%',
          marginLeft: -4
        }}
        classNames={'todo-list'}
      >
        {files.map((file) => {
          return (
            <CSSTransition
              key={file.checksum}
              timeout={200}
              classNames={'item'}
            >
              <SuggestedFile
                file={file}
                workspaceMnemonic={workspaceMnemonic}
              />
            </CSSTransition>
          );
        })}
      </TransitionGroup>
    </Box>
  );
};

export default SuggestedList;
