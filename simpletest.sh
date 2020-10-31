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
echo ""

echo "[TEST] test the example application:"
echo ""
cd ../
cd examples/example1
go run ../../runnr/runnr.go something -r
