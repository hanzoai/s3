/**
 * dl.hanzo.space — Download proxy for Hanzo S3 server binaries.
 *
 * Routes:
 *   /server/s3/release/linux-amd64/s3       → GitHub release asset s3-linux-amd64
 *   /server/s3/release/linux-arm64/s3       → GitHub release asset s3-linux-arm64
 *   /server/s3/release/darwin-amd64/s3      → GitHub release asset s3-darwin-amd64
 *   /server/s3/release/darwin-arm64/s3      → GitHub release asset s3-darwin-arm64
 *   /server/s3/release/windows-amd64/s3.exe → GitHub release asset s3-windows-amd64.exe
 */

const REPO = 'hanzoai/s3';
const CACHE_TTL = 3600; // 1 hour

const ASSET_MAP = {
  '/server/s3/release/linux-amd64/s3':       's3-linux-amd64',
  '/server/s3/release/linux-arm64/s3':       's3-linux-arm64',
  '/server/s3/release/darwin-amd64/s3':      's3-darwin-amd64',
  '/server/s3/release/darwin-arm64/s3':      's3-darwin-arm64',
  '/server/s3/release/windows-amd64/s3.exe': 's3-windows-amd64.exe',
};

export default {
  async fetch(request) {
    const url = new URL(request.url);
    const path = url.pathname;

    if (path === '/' || path === '') {
      return new Response('Hanzo S3 Downloads\nhttps://github.com/hanzoai/s3\n', {
        headers: { 'content-type': 'text/plain' },
      });
    }

    const assetName = ASSET_MAP[path];
    if (!assetName) {
      return new Response('Not found\n', { status: 404 });
    }

    // Get latest release
    const releaseResp = await fetch(
      `https://api.github.com/repos/${REPO}/releases/latest`,
      {
        headers: {
          'User-Agent': 'dl.hanzo.space',
          'Accept': 'application/vnd.github.v3+json',
        },
      }
    );

    if (!releaseResp.ok) {
      return new Response('Failed to fetch release info\n', { status: 502 });
    }

    const release = await releaseResp.json();
    const asset = release.assets?.find(a => a.name === assetName);

    if (!asset) {
      return new Response(`Asset ${assetName} not found in release ${release.tag_name}\n`, {
        status: 404,
      });
    }

    // Redirect to the download URL
    return Response.redirect(asset.browser_download_url, 302);
  },
};
