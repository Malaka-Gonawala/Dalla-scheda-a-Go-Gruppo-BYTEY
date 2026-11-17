param([string]$porta = "COM5")

$p = New-Object System.IO.Ports.SerialPort $porta,115200,None,8,One
$p.Open()
Write-Host "Aperta $porta. Scrivi un comando e premi Invio, 'esci' per chiudere."

$riga = ""
try {
  while ($true) {
    if ($p.BytesToRead -gt 0) { Write-Host -NoNewline $p.ReadExisting() }
    if ([Console]::KeyAvailable) {
      $k = [Console]::ReadKey($true)
      if ($k.Key -eq 'Enter') {
        Write-Host ""
        if ($riga -eq 'esci') { break }
        $p.Write($riga + "`n")
        $riga = ""
      } elseif ($k.Key -eq 'Backspace') {
        if ($riga.Length -gt 0) { $riga = $riga.Substring(0, $riga.Length - 1); Write-Host -NoNewline "`b `b" }
      } else {
        $riga += $k.KeyChar
        Write-Host -NoNewline $k.KeyChar
      }
    }
    Start-Sleep -Milliseconds 10
  }
} finally {
  $p.Close()
  Write-Host "Porta chiusa."
}
