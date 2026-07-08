import { describe, it, expect } from 'vitest';
import { cn } from '../utils';

describe('cn utility', () => {
  it('merges class names', () => {
    const result = cn('foo', 'bar');
    expect(result).toBe('foo bar');
  });

  it('handles conditional class names', () => {
    const result = cn('foo', false && 'hidden', 'bar');
    expect(result).toBe('foo bar');
  });

  it('deduplicates tailwind classes correctly', () => {
    // tailwind-merge should resolve conflicting utilities
    const result = cn('p-4', 'p-8');
    expect(result).toBe('p-8');
  });

  it('handles undefined and null gracefully', () => {
    const result = cn(undefined, null, 'foo');
    expect(result).toBe('foo');
  });

  it('returns empty string when no valid classes', () => {
    const result = cn(undefined, false, null);
    expect(result).toBe('');
  });
});
