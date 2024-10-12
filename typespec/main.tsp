import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";
import "./error.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

@service({
    title: "GlassCMS API",
})
namespace GlassCMSCore;

alias CreateItem = TypeSpec.Rest.Resource.ResourceCreate<Item, Error>;
alias GetItems = TypeSpec.Rest.Resource.ResourceRead<Item, Error>;
alias UpdateItem = TypeSpec.Rest.Resource.ResourceUpdate<Item, Error>;
alias ListItems = TypeSpec.Rest.Resource.ResourceList<Item, Error>;

interface Items {
    op create is CreateItem.create;
    op get is GetItems.get;
    op update is UpdateItem.update;

    @autoRoute
    @doc("Lists all instances of the resource.")
    @listsResource(Item)
    list(@query fields: string): Item[] | Error;
}

@doc("Item represents an individual content item.")
@resource("items")
model Item {
    @key
    @visibility("create", "read")
    name: string;
    display_name: string;
    content: string;
    create_time: utcDateTime;
    update_time: utcDateTime;
    @visibility("read")
    delete_time: utcDateTime;
    properties: Record<unknown>;
    metadata: Record<unknown>;
}
