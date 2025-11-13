import { uploadImage } from '@/utils';

describe('uploadImage', () => {
  const OLD_ENV = process.env;
  beforeEach(() => {
    jest.resetAllMocks();
    process.env = { ...OLD_ENV };
    process.env.NEXT_PUBLIC_CLOUDINARY_URL = 'https://api.cloudinary.com/v1_1/demo/upload';
    process.env.NEXT_PUBLIC_UPLOAD_PRESET = 'preset';
  });
  afterAll(() => {
    process.env = OLD_ENV;
  });

  it('returns secure_url on success', async () => {
    const mockJson = jest.fn().mockResolvedValue({ secure_url: 'https://cloudinary.com/image.jpg' });
    (global as any).fetch = jest.fn().mockResolvedValue({ ok: true, json: mockJson });
    const file = new File(['data'], 'pic.png', { type: 'image/png' });
    const url = await uploadImage(file);
    expect(url).toBe('https://cloudinary.com/image.jpg');
    expect(fetch).toHaveBeenCalled();
  });

  it('throws error body on failure', async () => {
    const errorBody = { error: 'bad request' };
    const mockJson = jest.fn().mockResolvedValue(errorBody);
    (global as any).fetch = jest.fn().mockResolvedValue({ ok: false, json: mockJson });
    const file = new File(['data'], 'pic.png', { type: 'image/png' });
    await expect(uploadImage(file)).rejects.toEqual(errorBody);
    expect(fetch).toHaveBeenCalled();
  });
});
