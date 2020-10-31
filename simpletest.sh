rm -rf tmp
mkdir tmp
cd tmp
echo ""
echo "[TEST] create new project:"
echo ""
go run ../runnr/runnr.go -n
echo ""
echo "[TEST] running help:"
echo ""
go run ../runnr/runnr.go -h
echo ""
echo "[TEST] running recompile and hello:"
echo ""
go run ../runnr/runnr.go hello -r
