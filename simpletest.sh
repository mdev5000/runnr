# NOTE: this was once used as a quick test, but tests have now been made actual test
#       so this is no longer need and just here for reference temporarily.
rm -rf _tmp
mkdir _tmp
cd _tmp
echo ""

echo "[TEST] create new project:"
echo ""
go run ../cmd/runnr/runnr.go g new
echo ""
echo "-------------------------------------------------------"
echo ""

echo "[TEST] running help:"
echo ""
go run ../cmd/runnr/runnr.go -h
echo ""
echo "-------------------------------------------------------"
echo ""

echo "[TEST] running recompile and hello:"
echo ""
go run ../cmd/runnr/runnr.go hello -r
echo ""
echo "-------------------------------------------------------"
echo ""

echo "[TEST] test the example application:"
echo ""
cd ../
cd examples/example1
go run ../../cmd/runnr/runnr.go something -r
echo ""
echo "-------------------------------------------------------"
echo ""

echo "[TEST] can run it statically"
echo ""
cd ../../
go run ./cmd/runnr/runnr.go g s run ./_tmp/internal/cmd/runnr_local/main.go ./_tmp/statictest -- hello -r
echo ""
echo "-------------------------------------------------------"
echo ""
