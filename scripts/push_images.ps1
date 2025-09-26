Param(
  [string]$User = "nallarmariano",
  [string]$StableTag = "v1.0",
  [string]$Latest = "latest",
  [string]$QaTag = "qa"
)

Write-Host "==> Publicando imágenes para usuario $User" -ForegroundColor Cyan

$images = @(
  @{ local="is3-backend:$StableTag"; remoteStable="$User/is3-backend:$StableTag"; remoteLatest="$User/is3-backend:$Latest" },
  @{ local="is3-frontend:$StableTag"; remoteStable="$User/is3-frontend:$StableTag"; remoteLatest="$User/is3-frontend:$Latest" },
  @{ local="is3-frontend:$QaTag"; remoteStable="$User/is3-frontend:$QaTag"; remoteLatest=$null }
)

foreach ($img in $images) {
  $local = $img.local
  if (-not (docker image inspect $local 2>$null | Out-Null)) {
    Write-Warning "Imagen local $local no existe. Saltando."
    continue
  }

  $remoteStable = $img.remoteStable
  Write-Host "Tag -> $local => $remoteStable" -ForegroundColor Yellow
  docker tag $local $remoteStable
  docker push $remoteStable

  if ($img.remoteLatest) {
    $remoteLatest = $img.remoteLatest
    Write-Host "Tag -> $local => $remoteLatest" -ForegroundColor Yellow
    docker tag $local $remoteLatest
    docker push $remoteLatest
  }
}

Write-Host "==> Finalizado" -ForegroundColor Green
