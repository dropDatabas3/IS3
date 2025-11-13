import React from 'react';
export const Swiper = ({ children }: { children: React.ReactNode }) => <div data-testid="swiper">{children}</div>;
export const SwiperSlide = ({ children }: { children: React.ReactNode }) => <div data-testid="slide">{children}</div>;
export default Swiper;
