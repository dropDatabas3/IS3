import { uiReducer } from '@/context/ui/uiReducer';

describe('uiReducer', () => {
  const initial = { isCreateModalOpen: false, isEdit: false };

  it('opens create modal', () => {
    const next = uiReducer(initial as any, { type: '[Ui] - Open Create Modal' });
    expect(next.isCreateModalOpen).toBe(true);
    expect(next.isEdit).toBe(false);
  });
  it('opens edit modal', () => {
    const next = uiReducer(initial as any, { type: '[Ui] - Open Edit Modal' });
    expect(next.isCreateModalOpen).toBe(true);
    expect(next.isEdit).toBe(true);
  });
  it('closes create modal', () => {
    const open = { isCreateModalOpen: true, isEdit: false };
    const next = uiReducer(open as any, { type: '[Ui] - Close Create Modal' });
    expect(next.isCreateModalOpen).toBe(false);
  });
  it('closes edit modal', () => {
    const openEdit = { isCreateModalOpen: true, isEdit: true };
    const next = uiReducer(openEdit as any, { type: '[Ui] - Close Edit Modal' });
    expect(next.isEdit).toBe(false);
    expect(next.isCreateModalOpen).toBe(false);
  });
  it('default returns same state', () => {
    // @ts-expect-error unknown
    expect(uiReducer(initial as any, { type: 'UNKNOWN' })).toEqual(initial);
  });
});
