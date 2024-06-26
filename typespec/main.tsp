import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

using TypeSpec.Http;
using TypeSpec.Rest;

@service({
    title: "GlassCMS API",
})
namespace GlassCMS;

@route("/items")
namespace Items {
    @get op list(): Item[];
    @delete op delete(id: int64): {};
    @post op create(...Item): {};
}

@@listsResource(Items.list, Item);
@@createsResource(Items.create, Item);

@doc("Item represents an individual content item.")
model Item {
    id: string;
    name: string;
    display_name: string;
    path: string;
    content: string;
    create_time: utcDateTime;
    update_time: utcDateTime;
    properties: Record<unknown>;
}
