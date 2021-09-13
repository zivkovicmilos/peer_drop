import { Box, Divider, Typography } from '@material-ui/core';
import React, { FC, useEffect, useRef, useState } from 'react';
import { ISectionTitleProps } from './sectionTitle.types';

const SectionTitle: FC<ISectionTitleProps> = (props) => {
  const {
    title,
    titleClass = 'sectionTitle',
    dividerClass = 'sectionDivider'
  } = props;

  const [titleWidth, setTitleWidth] = useState<string>('0');

  const ref = useRef<HTMLHeadingElement>(null);

  const updateTitleWidth = () => {
    const width = ref.current ? ref.current.offsetWidth : 0;

    setTitleWidth(width * 0.9 + 'px');
  };

  useEffect(() => {
    updateTitleWidth();
  }, [ref.current, title]);

  useEffect(() => {
    window.addEventListener('resize', updateTitleWidth);
    return () => window.removeEventListener('resize', updateTitleWidth);
  }, []);

  return (
    <Box display={'block'} flexDirection={'column'}>
      <Typography className={titleClass} component={'span'} ref={ref}>
        {title}
      </Typography>
      <Divider
        style={{
          width: titleWidth
        }}
        classes={{
          root: dividerClass
        }}
      />
    </Box>
  );
};

export default SectionTitle;